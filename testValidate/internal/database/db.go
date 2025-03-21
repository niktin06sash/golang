package database

import (
	"database/sql"
	"fmt"
	"os"
	"testValidate/internal/erro"
)

type DBInterface interface {
	Open(driverName, dataSourceName string) (*sql.DB, error)
}

type DBObject struct{}

func (d DBObject) Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

type EnvInterface interface {
	GetEnv(key string) (string, error)
}
type EnvObject struct{}

func (ev EnvObject) GetEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", erro.ErrorGetEnv
	}
	return value, nil
}
func ConnectToDb(dbInterface DBInterface, env EnvInterface) (*sql.DB, error) {
	dbUser, err := env.GetEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPass, err := env.GetEnv("DB_PASS")
	if err != nil {
		return nil, err
	}
	dbName, err := env.GetEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", dbUser, dbPass, dbName)

	db, err := dbInterface.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
