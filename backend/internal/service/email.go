package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/user/votex-template/backend/internal/config"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (s *EmailService) SendPasswordResetEmail(email, token string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		Hello,
		
		You have requested a password reset for your account.
		
		Click the following link to reset your password:
		%s/auth/reset-password?token=%s
		
		This link will expire in %d hours.
		
		If you did not request this reset, please ignore this email.
		
		Best regards,
		The Vortex Team
	`, s.config.AppURL, token, s.config.PasswordResetTokenExpiry)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) SendWelcomeEmail(email, username string) error {
	subject := "Welcome to Vortex!"
	body := fmt.Sprintf(`
		Hello %s,
		
		Welcome to Vortex! Your account has been successfully created.
		
		You can now log in to your account and start using our services.
		
		If you have any questions, please don't hesitate to contact us.
		
		Best regards,
		The Vortex Team
	`, username)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	// If SMTP is not configured, just log the email (for development)
	if s.config.SMTPHost == "" {
		fmt.Printf("Email would be sent to %s:\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	// Create message
	message := fmt.Sprintf("From: %s\r\n", s.config.SMTPFrom)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n"
	message += body

	// Setup authentication
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	// Connect to server
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	var err error
	if s.config.SMTPTLS {
		// TLS connection
		tlsConfig := &tls.Config{
			ServerName: s.config.SMTPHost,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err = client.Mail(s.config.SMTPFrom); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}
		defer writer.Close()

		_, err = writer.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	} else {
		// Non-TLS connection
		err = smtp.SendMail(addr, auth, s.config.SMTPFrom, []string{to}, []byte(message))
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}
