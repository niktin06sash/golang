package server

import (
	"context"
	"database/sql"
	"net/http"
	"reflect"
	"testing"

	"testValidate/internal/config"
	"testValidate/internal/database"
	"testValidate/internal/person"

	"github.com/golang/mock/gomock" // Используем gomock
	"github.com/stretchr/testify/assert"
)

// Определяем интерфейс для PersonService (если его еще нет)
type MockPersonService interface {
	Registration(ctx context.Context, newperk person.Person) error
	// Другие методы PersonService
}

// MockDBRepository - Mock для DBRepository
type MockDBRepository struct {
	database.PersonRepository // Embed the interface to satisfy it
	Ctrl                      *gomock.Controller
}

type MockDBRepositoryMockRecorder struct {
	mock *MockDBRepository
}

// NewMockDBRepository creates a new mock instance
func NewMockDBRepository(ctrl *gomock.Controller) *MockDBRepository {
	mock := &MockDBRepository{Ctrl: ctrl}

	return mock
}

// Get mocks base method
func (m *MockDBRepository) Get(ctx context.Context) (*sql.Rows, error) {
	ret := m.Ctrl.Call(m, "Get", ctx)
	r0, _ := ret[0].(*sql.Rows)
	r1, _ := ret[1].(error)
	return r0, r1
}

// Mock DatabaseChecker
type MockDatabaseChecker struct {
	PingContextFunc func(ctx context.Context) error
}

func (m *MockDatabaseChecker) PingContext(ctx context.Context) error {
	if m.PingContextFunc != nil {
		return m.PingContextFunc(ctx)
	}
	return nil
}

func TestNewServer(t *testing.T) {
	// Arrange
	// Создаем mock Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock PersonRepository
	mockRepo := NewMockDBRepository(ctrl)

	// Создаем mock sql.DB
	db := &sql.DB{}

	// Создаем конфигурацию
	cfg := &config.Config{Port: ":8080"}

	// Act
	server := NewServer(cfg, db)

	// Assert
	// Проверяем, что поля структуры Server были инициализированы правильно
	assert.NotNil(t, server)
	assert.IsType(t, &person.PersonService{}, server.PersonService)
	assert.IsType(t, &http.Server{}, server.HttpServer)
	assert.Equal(t, cfg.Port, server.HttpServer.Addr)

	// Проверяем, что PersonService был создан с правильным репозиторием
	personService := server.PersonService
	// Используем reflection для получения значения поля Repo структуры PersonService
	repoValue := reflect.ValueOf(personService).Elem().FieldByName("Repo")

	// Сравниваем значения поля Repo с mockRepo
	assert.Equal(t, mockRepo, repoValue.Interface())
}
