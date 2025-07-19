package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/user/votex-template/backend/internal/config"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	config   *config.Config
}

func NewRateLimiter(cfg *config.Config) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		config:   cfg,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		window := time.Duration(rl.config.RateLimitRequests) * time.Minute

		for ip, requests := range rl.requests {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < window {
					validRequests = append(validRequests, reqTime)
				}
			}
			if len(validRequests) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		ip := getClientIP(r)

		rl.mutex.Lock()
		now := time.Now()
		window := time.Duration(rl.config.RateLimitRequests) * time.Minute

		// Get existing requests for this IP
		requests, exists := rl.requests[ip]
		if !exists {
			requests = []time.Time{}
		}

		// Remove old requests outside the window
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < window {
				validRequests = append(validRequests, reqTime)
			}
		}

		// Check if rate limit exceeded
		if len(validRequests) >= rl.config.RateLimitBurst {
			rl.mutex.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-RateLimit-Limit", string(rune(rl.config.RateLimitBurst)))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", time.Now().Add(window).Format(time.RFC3339))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Rate limit exceeded. Please try again later."}`))
			return
		}

		// Add current request
		validRequests = append(validRequests, now)
		rl.requests[ip] = validRequests

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", string(rune(rl.config.RateLimitBurst)))
		w.Header().Set("X-RateLimit-Remaining", string(rune(rl.config.RateLimitBurst-len(validRequests))))
		w.Header().Set("X-RateLimit-Reset", time.Now().Add(window).Format(time.RFC3339))

		rl.mutex.Unlock()

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	// Check for forwarded headers (for proxy/load balancer scenarios)
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}

	// Fallback to remote address
	return r.RemoteAddr
}
