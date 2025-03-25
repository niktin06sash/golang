package repository

import (
	"database/sql"
	"microservicesProject/auth_service/internal/model"
)

type Authorization interface {
	CreateUser(user *model.Person) (int, error)
	GetUser(useremail, password string) (*model.Person, error)
}
type AuthRepository struct {
	Authorization
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		Authorization: NewAuthPostgres(db),
	}
}
