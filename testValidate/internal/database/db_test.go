package database_test

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"testValidate/internal/database"
	"testValidate/internal/erro"

	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	OpenFunc func(driverName, dataSourceName string) (*sql.DB, error)
}

func (m MockDB) Open(driverName string, dataSourceName string) (*sql.DB, error) {
	return m.OpenFunc(driverName, dataSourceName)
}

type MockEnv struct {
	GetEnvFunc func(key string) (string, error)
}

func (m MockEnv) GetEnv(key string) (string, error) {
	return m.GetEnvFunc(key)
}

func TestGetEnv_Success(t *testing.T) {
	key := "TEST_ENV_VAR"
	value := "test_value"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	envObject := database.EnvObject{}

	result, err := envObject.GetEnv(key)

	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestGetEnv_Error(t *testing.T) {
	key := "NON_EXISTENT_ENV_VAR"

	envObject := database.EnvObject{}

	result, err := envObject.GetEnv(key)

	assert.Error(t, err)
	assert.ErrorIs(t, err, erro.ErrorGetEnv)
	assert.Equal(t, "", result)
}
func TestConnectToDb_ErrorSqlOpen(t *testing.T) {
	mockDB := MockDB{
		OpenFunc: func(driverName string, dataSourceName string) (*sql.DB, error) {
			return nil, errors.New("failed to connect to database")
		},
	}

	mockEnv := MockEnv{
		GetEnvFunc: func(key string) (string, error) {
			return "value", nil // Provide a value so ConnectToDb doesn't fail before Open
		},
	}

	db, err := database.ConnectToDb(mockDB, mockEnv)

	assert.Error(t, err)
	assert.Nil(t, db)
}
func TestConnectToDb_Success(t *testing.T) {

	mockEnv := MockEnv{
		GetEnvFunc: func(key string) (string, error) {
			switch key {
			case "DB_USER":
				return "user", nil
			case "DB_PASS":
				return "password", nil
			case "DB_NAME":
				return "dbname", nil
			default:
				return "", errors.New("unexpected key")
			}
		},
	}
	dbObject := database.DBObject{} // Use DBObject, not MockDB

	db, err := database.ConnectToDb(dbObject, mockEnv)

	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestConnectToDb_ErrorGetEnv_DBUser(t *testing.T) {

	mockEnv := MockEnv{
		GetEnvFunc: func(key string) (string, error) {
			if key == "DB_USER" {
				return "", erro.ErrorGetEnv
			}
			return "value", nil
		},
	}

	dbObject := database.DBObject{} // Use DBObject, not MockDB

	db, err := database.ConnectToDb(dbObject, mockEnv)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestConnectToDb_ErrorGetEnv_DBPass(t *testing.T) {

	mockEnv := MockEnv{
		GetEnvFunc: func(key string) (string, error) {
			if key == "DB_PASS" {
				return "", erro.ErrorGetEnv
			}
			return "value", nil
		},
	}

	dbObject := database.DBObject{} // Use DBObject, not MockDB

	db, err := database.ConnectToDb(dbObject, mockEnv)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestConnectToDb_ErrorGetEnv_DBName(t *testing.T) {

	mockEnv := MockEnv{
		GetEnvFunc: func(key string) (string, error) {
			if key == "DB_NAME" {
				return "", erro.ErrorGetEnv
			}
			return "value", nil
		},
	}

	dbObject := database.DBObject{} // Use DBObject, not MockDB

	db, err := database.ConnectToDb(dbObject, mockEnv)

	assert.Error(t, err)
	assert.Nil(t, db)
}
