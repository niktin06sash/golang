package repository

import (
	"context"
	"database/sql"
	"microservicesProject/auth_service/internal/model"
)

type Authorization interface {
	CreateUser(ctx context.Context, user *model.Person) *AuthenticationRepositoryResponse
	GetUser(ctx context.Context, useremail, password string) *AuthenticationRepositoryResponse
}
type Repository struct {
	Authorization
}

func NewAuthRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
