package service

import (
	"context"
	"microservicesProject/auth_service/internal/model"
	"microservicesProject/auth_service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type Authorization interface {
	Registrate(user *model.Person, ctx context.Context) *AuthenticationServiceResponse
	Authorizate(user *model.Person, ctx context.Context) *AuthenticationServiceResponse
	GenerateSession(userId uuid.UUID) (string, time.Time)
	CheckSession(sessionID string) *AuthenticationServiceResponse
	//DeleteSession(token string) error
}
type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {

	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
