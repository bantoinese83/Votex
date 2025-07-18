package service

import "github.com/user/votex-template/backend/internal/store"

type Service struct {
	Store *store.Store
}

func New(s *store.Store) *Service {
	return &Service{Store: s}
}
