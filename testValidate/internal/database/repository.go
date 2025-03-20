package database

import (
	"context"
	"database/sql"

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
	Add(Name string, Email string, Password string, ctx context.Context) error
	CheckUniqueEmail(Email string, ctx context.Context) (bool, error)
	CompareEmail(Email string, ctx context.Context) (bool, error)
	ComparePassword(Email string, password string, ctx context.Context) (bool, error)
}

func (repos *DBRepository) Add(Name string, Email string, Password string, ctx context.Context) error {
	_, err := repos.DataBase.ExecContext(ctx, "INSERT INTO UserZ (username, useremail, userpassword) values ($1, $2, $3)", Name, Email, Password)

	if err != nil {
		return err
	}
	return nil
}

func (repos *DBRepository) CheckUniqueEmail(Email string, ctx context.Context) (bool, error) {
	rows, err := repos.DataBase.QueryContext(ctx, "SELECT * from userZ where useremail = $1", Email)
	if err != nil {

		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return true, nil
	}
	return false, nil
}
func (repos *DBRepository) CompareEmail(Email string, ctx context.Context) (bool, error) {
	rows, err := repos.DataBase.QueryContext(ctx, "SELECT * from userZ where useremail = $1", Email)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, err
	}
	return true, nil
}
func (repos *DBRepository) ComparePassword(Email string, password string, ctx context.Context) (bool, error) {
	rows, err := repos.DataBase.QueryContext(ctx, "SELECT userpassword from userZ where useremail = $1", Email)
	if err != nil {
		return false, err
	}
	var hashpass string

	defer rows.Close()
	if !rows.Next() {
		return false, nil
	}
	err = rows.Scan(&hashpass)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password))
	if err == nil {
		return true, nil
	}
	return false, err
}
