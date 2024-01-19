package service

import (
	"testApp/internal/repository"
	"testApp/internal/service/users"
)

type Service struct {
	User users.UserService
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		User: users.NewUserService(r.Users),
	}
}
