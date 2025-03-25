package repository

import (
	"database/sql"
	"microservicesProject/auth_service/internal/model"
)

type AuthPostgres struct {
	Db *sql.DB
}

func (ap *AuthPostgres) CreateUser(user *model.Person) (int, error) {
	return 0, nil
}
func (ap *AuthPostgres) GetUser(useremail, password string) (*model.Person, error) {
	return nil, nil
}
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{Db: db}
}
