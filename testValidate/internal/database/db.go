package database

import (
	"database/sql"
	"fmt"
	"os"
	"testValidate/internal/erro"
)

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{
		DataBase: db,
	}
}
func getEnv(key string) (string, error) {

	value := os.Getenv(key)

	if value == "" {
		return "", erro.ErrorGetEnv
	}
	return value, nil
}
func ConnectToDb() (*sql.DB, error) {
	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPass, err := getEnv("DB_PASS")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}
	sslMode := "disable"
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", dbUser, dbName, dbPass, sslMode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, erro.ErrorDBConnect
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, erro.ErrorDBPing
	}
	return db, nil
}
