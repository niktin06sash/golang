package service

import (
	"microservicesProject/auth_service/internal/model"
	"microservicesProject/auth_service/internal/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
func (as *AuthService) CreateUser(user model.Person) (int, error) {
	return 0, nil
}
func (as *AuthService) GenerateSession(username, password string) (string, error) {
	return "", nil
}
func (as *AuthService) CheckSession(token string) (int, error) {
	return 0, nil
}
func (as *AuthService) DeleteSession(token string) error {
	return nil
}
