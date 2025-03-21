package database

import (
	"context"
	"database/sql"
	"errors"
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
	Add(UserId uuid.UUID, Name string, Email string, Password string, ctx context.Context) error
	AuthenticateUser(Email string, password string, ctx context.Context) (bool, error)
}

func (repos *DBRepository) Add(UserId uuid.UUID, Name string, Email string, Password string, ctx context.Context) error {

	result, err := repos.DataBase.ExecContext(ctx, "INSERT INTO UserZ (userid, username, useremail, userpassword) values ($1, $2, $3, $4) ON CONFLICT (useremail) DO NOTHING;", UserId, Name, Email, Password)
	if err != nil {

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return err
	}

	if rowsAffected == 0 {

		return erro.ErrorUniqueEmail
	}

	return nil

}

func (repos *DBRepository) AuthenticateUser(Email string, password string, ctx context.Context) (bool, error) {
	var hashpass string
	err := repos.DataBase.QueryRowContext(ctx, "SELECT userpassword FROM userZ WHERE useremail = $1", Email).Scan(&hashpass)

	if err == sql.ErrNoRows {

		return false, erro.ErrorEmailNotRegister
	}
	if err != nil {

		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {

		return false, erro.ErrorInvalidPerson
	}

	return false, err
}
