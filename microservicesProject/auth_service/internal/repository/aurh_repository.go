package repository

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthPostgres struct {
	Db *sql.DB
}
type AuthenticationRepositoryResponse struct {
	Success bool
	UserId  uuid.UUID
	Errors  error
}

func (repoap *AuthPostgres) CreateUser(ctx context.Context, user *model.Person) *AuthenticationRepositoryResponse {
	result, err := repoap.Db.ExecContext(ctx, "INSERT INTO UserZ (userid, username, useremail, userpassword) values ($1, $2, $3, $4) ON CONFLICT (useremail) DO NOTHING;", user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return &AuthenticationRepositoryResponse{Success: false, Errors: err}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return &AuthenticationRepositoryResponse{Success: false, Errors: err}
	}

	if rowsAffected == 0 {
		return &AuthenticationRepositoryResponse{Success: false, Errors: erro.ErrorUniqueEmail}
	}

	return &AuthenticationRepositoryResponse{Success: true, Errors: nil, UserId: user.Id}
}
func (repoap *AuthPostgres) GetUser(ctx context.Context, useremail, userpassword string) *AuthenticationRepositoryResponse {
	var hashpass string
	var userId uuid.UUID
	err := repoap.Db.QueryRowContext(ctx, "SELECT userid, userpassword FROM userZ WHERE useremail = $1", useremail).Scan(&userId, &hashpass)

	if err == sql.ErrNoRows {
		return &AuthenticationRepositoryResponse{UserId: uuid.Nil, Success: false, Errors: erro.ErrorEmailNotRegister}
	}
	if err != nil {
		return &AuthenticationRepositoryResponse{UserId: uuid.Nil, Success: false, Errors: err}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(userpassword))
	if err != nil {

		return &AuthenticationRepositoryResponse{UserId: uuid.Nil, Success: false, Errors: erro.ErrorInvalidPassword}
	}

	return &AuthenticationRepositoryResponse{UserId: userId, Success: true, Errors: nil}
}
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{Db: db}
}
