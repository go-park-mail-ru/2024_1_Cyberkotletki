// Code generated by MockGen. DO NOT EDIT.
// Source: content.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	gomock "github.com/golang/mock/gomock"
)

// MockContent is a mock of Content interface.
type MockContent struct {
	ctrl     *gomock.Controller
	recorder *MockContentMockRecorder
}

// MockContentMockRecorder is the mock recorder for MockContent.
type MockContentMockRecorder struct {
	mock *MockContent
}

// NewMockContent creates a new mock instance.
func NewMockContent(ctrl *gomock.Controller) *MockContent {
	mock := &MockContent{ctrl: ctrl}
	mock.recorder = &MockContentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContent) EXPECT() *MockContentMockRecorder {
	return m.recorder
}

// GetContentByID mocks base method.
func (m *MockContent) GetContentByID(id int) (*dto.Content, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentByID", id)
	ret0, _ := ret[0].(*dto.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContentByID indicates an expected call of GetContentByID.
func (mr *MockContentMockRecorder) GetContentByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentByID", reflect.TypeOf((*MockContent)(nil).GetContentByID), id)
}

// GetPersonByID mocks base method.
func (m *MockContent) GetPersonByID(id int) (*dto.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonByID", id)
	ret0, _ := ret[0].(*dto.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersonByID indicates an expected call of GetPersonByID.
func (mr *MockContentMockRecorder) GetPersonByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonByID", reflect.TypeOf((*MockContent)(nil).GetPersonByID), id)
}
