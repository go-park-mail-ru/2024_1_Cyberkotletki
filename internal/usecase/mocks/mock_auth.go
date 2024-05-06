// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go
//
// Generated by this command:
//
//	mockgen -source=auth.go -destination=mocks/mock_auth.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockAuth) CreateSession(userID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockAuthMockRecorder) CreateSession(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockAuth)(nil).CreateSession), userID)
}

// GetUserIDBySession mocks base method.
func (m *MockAuth) GetUserIDBySession(session string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDBySession", session)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDBySession indicates an expected call of GetUserIDBySession.
func (mr *MockAuthMockRecorder) GetUserIDBySession(session any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDBySession", reflect.TypeOf((*MockAuth)(nil).GetUserIDBySession), session)
}

// Logout mocks base method.
func (m *MockAuth) Logout(session string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", session)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthMockRecorder) Logout(session any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuth)(nil).Logout), session)
}

// LogoutAll mocks base method.
func (m *MockAuth) LogoutAll(userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutAll", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogoutAll indicates an expected call of LogoutAll.
func (mr *MockAuthMockRecorder) LogoutAll(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutAll", reflect.TypeOf((*MockAuth)(nil).LogoutAll), userID)
}
