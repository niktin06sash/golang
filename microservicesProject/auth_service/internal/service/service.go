package service

import (
	"microservicesProject/auth_service/internal/model"
	"microservicesProject/auth_service/internal/repository"
)

type Authorization interface {
	CreateUser(user model.Person) (int, error)
	GenerateSession(username, password string) (string, error)
	CheckSession(token string) (int, error)
	DeleteSession(token string) error
}
type Service struct {
	Authorization
}

func NewService(repos *repository.AuthRepository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
