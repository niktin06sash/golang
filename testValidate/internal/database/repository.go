package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type PersonRepository interface {
	Add(Name string, Email string, ctx context.Context) error
	Get(ctx context.Context) (*sql.Rows, error)
}
type DBRepository struct {
	DataBase *sql.DB
}

func (repos *DBRepository) Add(Name string, Email string, ctx context.Context) error {
	_, err := repos.DataBase.ExecContext(ctx, "INSERT INTO testGo (username, useremail) values ($1, $2)", Name, Email)

	if err != nil {
		return err
	}
	return nil
}
func (repos *DBRepository) Get(ctx context.Context) (*sql.Rows, error) {

	rows, err := repos.DataBase.QueryContext(ctx, "SELECT * FROM testGo")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
