package database

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"testValidate/internal/erro" // Замените на правильный путь к вашему пакету erro

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestDBRepository_Add_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	name := "Test User"
	email := "test@example.com"
	password := "password"

	// Генерируем UUID перед вызовом Add
	userID := uuid.New()

	mock.ExpectExec("INSERT INTO UserZ").
		WithArgs(
			userID, // Проверяем конкретное значение UUID
			name,
			email,
			password,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add(userID, name, email, password, context.Background()) // Передаем UUID
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBRepository_Add_Conflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	name := "Test User"
	email := "test@example.com"
	password := "password"

	// Генерируем UUID перед вызовом Add
	userID := uuid.New()

	mock.ExpectExec("INSERT INTO UserZ").
		WithArgs(
			userID, // Проверяем конкретное значение UUID
			name,
			email,
			password,
		).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

	err = repo.Add(userID, name, email, password, context.Background()) // Передаем UUID
	assert.Error(t, err)
	assert.Equal(t, erro.ErrorUniqueEmail, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBRepository_Add_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	name := "Test User"
	email := "test@example.com"
	password := "password"

	// Генерируем UUID перед вызовом Add
	userID := uuid.New()

	mock.ExpectExec("INSERT INTO UserZ").
		WithArgs(
			userID, // Проверяем конкретное значение UUID
			name,
			email,
			password,
		).
		WillReturnError(errors.New("database error"))

	err = repo.Add(userID, name, email, password, context.Background()) // Передаем UUID
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBRepository_AuthenticateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	email := "test@example.com"
	password := "password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT userpassword FROM userZ WHERE useremail = $1")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"userpassword"}).AddRow(hashedPassword))

	success, err := repo.AuthenticateUser(email, password, context.Background())
	assert.NoError(t, err)
	assert.True(t, success)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBRepository_AuthenticateUser_EmailNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	email := "test@example.com"
	password := "password"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT userpassword FROM userZ WHERE useremail = $1")).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	success, err := repo.AuthenticateUser(email, password, context.Background())
	assert.Error(t, err)
	assert.Equal(t, erro.ErrorEmailNotRegister, err)
	assert.False(t, success)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBRepository_AuthenticateUser_InvalidPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	email := "test@example.com"
	password := "wrongpassword"
	correctPassword := "password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT userpassword FROM userZ WHERE useremail = $1")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"userpassword"}).AddRow(hashedPassword))

	success, err := repo.AuthenticateUser(email, password, context.Background())
	assert.Error(t, err)
	assert.Equal(t, erro.ErrorInvalidPerson, err)
	assert.False(t, success)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBRepository_AuthenticateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	email := "test@example.com"
	password := "password"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT userpassword FROM userZ WHERE useremail = $1")).
		WithArgs(email).
		WillReturnError(errors.New("database error"))

	success, err := repo.AuthenticateUser(email, password, context.Background())
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.False(t, success)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
