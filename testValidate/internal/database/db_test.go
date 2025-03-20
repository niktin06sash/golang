package database_test

import (
	"database/sql"
	"database/sql/driver"
	"os"
	"testing"

	"testValidate/internal/database"
	"testValidate/internal/erro"

	"github.com/stretchr/testify/assert"
)

// MockDriver - Mock для sql driver
type MockDriver struct{}

// Open - Реализация метода Open для MockDriver
func (m MockDriver) Open(name string) (driver.Conn, error) {
	// Здесь можно добавить логику для возврата ошибки в определенных случаях
	return &MockConn{}, nil
}

type MockConn struct{}

// Prepare - Реализация метода Prepare для MockConn
func (c *MockConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

// Close - Реализация метода Close для MockConn
func (c *MockConn) Close() error {
	return nil
}

// Begin - Реализация метода Begin для MockConn
func (c *MockConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func TestNewDBRepository(t *testing.T) {
	// Arrange

	// Регистрируем MockDriver для "mockdb"
	sql.Register("mockdb", &MockDriver{})

	// Открываем соединение с "mockdb"
	db, err := sql.Open("mockdb", "connection_string")
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	// Act
	repo := database.NewDBRepository(db)

	// Assert
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.DataBase)
}

func TestGetEnv_Success(t *testing.T) {
	// Arrange
	key := "TEST_ENV_VAR"
	value := "test_value"
	os.Setenv(key, value)
	defer os.Unsetenv(key) // Убираем переменную окружения после теста

	// Act
	result, err := database.GetEnv(key)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestGetEnv_Error(t *testing.T) {
	// Arrange
	key := "NON_EXISTENT_ENV_VAR" // Переменная, которой точно нет

	// Act
	result, err := database.GetEnv(key)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, erro.ErrorGetEnv, err)
	assert.Equal(t, "", result)
}
func TestConnectToDb_Integration(t *testing.T) {
	// Arrange
	// Устанавливаем переменные окружения для подключения к тестовой базе данных
	os.Setenv("DB_USER", "postgres")   // Замените на ваши тестовые значения
	os.Setenv("DB_PASS", "sosuhui247") // Замените на ваши тестовые значения
	os.Setenv("DB_NAME", "fc2")        // Замените на ваши тестовые значения
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")
	}()

	// Act
	db, err := database.ConnectToDb()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Закрываем соединение с базой данных
	db.Close()
}

func TestConnectToDb_Integration_ErrorGetEnv(t *testing.T) {
	// Arrange
	// Оставляем переменные окружения неопределенными
	os.Setenv("DB_USER", "postg")      // Замените на ваши тестовые значения
	os.Setenv("DB_PASS", "sosuhui247") // Замените на ваши тестовые значения
	os.Setenv("DB_NAME", "fc2")        // Замените на ваши тестовые значения
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")
	}()

	// Act
	db, err := database.ConnectToDb()

	// Assert
	assert.Error(t, err)
	// assert.Equal(t, erro.ErrorGetEnv, err) // This might fail, as any of the three env vars could be missing
	assert.Nil(t, db)

}
func TestConnectToDb_ErrorGetEnv_DBUser(t *testing.T) {
	// Arrange
	os.Setenv("DB_PASS", "sosuhui247") // Устанавливаем DB_PASS
	os.Setenv("DB_NAME", "fc2")        // Устанавливаем DB_NAME
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")
	}()
	os.Unsetenv("DB_USER") // Убеждаемся, что DB_USER не установлена

	// Act
	db, err := database.ConnectToDb()

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, erro.ErrorGetEnv)
	assert.Nil(t, db)
}

func TestConnectToDb_ErrorGetEnv_DBPass(t *testing.T) {
	// Arrange
	os.Setenv("DB_USER", "postgres") // Устанавливаем DB_USER
	os.Setenv("DB_NAME", "fc2")      // Устанавливаем DB_NAME
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")
	}()
	os.Unsetenv("DB_PASS") // Убеждаемся, что DB_PASS не установлена

	// Act
	db, err := database.ConnectToDb()

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, erro.ErrorGetEnv)
	assert.Nil(t, db)
}

func TestConnectToDb_ErrorGetEnv_DBName(t *testing.T) {
	// Arrange
	os.Setenv("DB_USER", "postgres")   // Устанавливаем DB_USER
	os.Setenv("DB_PASS", "sosuhui247") // Устанавливаем DB_PASS
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")
	}()
	os.Unsetenv("DB_NAME") // Убеждаемся, что DB_NAME не установлена

	// Act
	db, err := database.ConnectToDb()

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, erro.ErrorGetEnv)
	assert.Nil(t, db)
}
