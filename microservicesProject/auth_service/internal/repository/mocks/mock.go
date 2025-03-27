// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "auth_service/internal/model"
	repository "auth_service/internal/repository"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthorizationRepos is a mock of AuthorizationRepos interface.
type MockAuthorizationRepos struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationReposMockRecorder
}

// MockAuthorizationReposMockRecorder is the mock recorder for MockAuthorizationRepos.
type MockAuthorizationReposMockRecorder struct {
	mock *MockAuthorizationRepos
}

// NewMockAuthorizationRepos creates a new mock instance.
func NewMockAuthorizationRepos(ctrl *gomock.Controller) *MockAuthorizationRepos {
	mock := &MockAuthorizationRepos{ctrl: ctrl}
	mock.recorder = &MockAuthorizationReposMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizationRepos) EXPECT() *MockAuthorizationReposMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthorizationRepos) CreateUser(ctx context.Context, user *model.Person) *repository.AuthenticationRepositoryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(*repository.AuthenticationRepositoryResponse)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationReposMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorizationRepos)(nil).CreateUser), ctx, user)
}

// GetUser mocks base method.
func (m *MockAuthorizationRepos) GetUser(ctx context.Context, useremail, password string) *repository.AuthenticationRepositoryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, useremail, password)
	ret0, _ := ret[0].(*repository.AuthenticationRepositoryResponse)
	return ret0
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAuthorizationReposMockRecorder) GetUser(ctx, useremail, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuthorizationRepos)(nil).GetUser), ctx, useremail, password)
}
