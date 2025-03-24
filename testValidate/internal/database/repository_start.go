package database

import (
	"context"
	"database/sql"
	"testValidate/internal/erro"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{
		DataBase: db,
	}
}

type DBRepository struct {
	DataBase *sql.DB
}
type PersonRepository interface {
	Add(UserId uuid.UUID, Name string, Email string, Password string, ctx context.Context) *AddResult
	AuthenticateUser(Email string, password string, ctx context.Context) *AuthenticationResult
}

type AddResult struct {
	Success bool
	Error   error
}

func (repos *DBRepository) Add(UserId uuid.UUID, Name string, Email string, Password string, ctx context.Context) *AddResult {
	result, err := repos.DataBase.ExecContext(ctx, "INSERT INTO UserZ (userid, username, useremail, userpassword) values ($1, $2, $3, $4) ON CONFLICT (useremail) DO NOTHING;", UserId, Name, Email, Password)
	if err != nil {
		return &AddResult{Success: false, Error: err}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return &AddResult{Success: false, Error: err}
	}

	if rowsAffected == 0 {
		return &AddResult{Success: false, Error: erro.ErrorUniqueEmail}
	}

	return &AddResult{Success: true, Error: nil}
}

type AuthenticationResult struct {
	UserID  uuid.UUID
	Success bool
	Error   error
}

func (repos *DBRepository) AuthenticateUser(Email string, password string, ctx context.Context) *AuthenticationResult {
	var hashpass string
	var userId uuid.UUID
	err := repos.DataBase.QueryRowContext(ctx, "SELECT userid, userpassword FROM userZ WHERE useremail = $1", Email).Scan(&userId, &hashpass)

	if err == sql.ErrNoRows {
		return &AuthenticationResult{UserID: uuid.Nil, Success: false, Error: erro.ErrorEmailNotRegister}
	}
	if err != nil {
		return &AuthenticationResult{UserID: uuid.Nil, Success: false, Error: err}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password))
	if err != nil {

		return &AuthenticationResult{UserID: uuid.Nil, Success: false, Error: erro.ErrorInvalidPerson}
	}

	return &AuthenticationResult{UserID: userId, Success: true, Error: nil}
}
