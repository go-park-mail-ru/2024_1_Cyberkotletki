// Code generated by MockGen. DO NOT EDIT.
// Source: compilation.go
//
// Generated by this command:
//
//	mockgen -source=compilation.go -destination=mocks/mock_compilation.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	entity "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockCompilation is a mock of Compilation interface.
type MockCompilation struct {
	ctrl     *gomock.Controller
	recorder *MockCompilationMockRecorder
}

// MockCompilationMockRecorder is the mock recorder for MockCompilation.
type MockCompilationMockRecorder struct {
	mock *MockCompilation
}

// NewMockCompilation creates a new mock instance.
func NewMockCompilation(ctrl *gomock.Controller) *MockCompilation {
	mock := &MockCompilation{ctrl: ctrl}
	mock.recorder = &MockCompilationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompilation) EXPECT() *MockCompilationMockRecorder {
	return m.recorder
}

// GetAllCompilationTypes mocks base method.
func (m *MockCompilation) GetAllCompilationTypes() ([]entity.CompilationType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCompilationTypes")
	ret0, _ := ret[0].([]entity.CompilationType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCompilationTypes indicates an expected call of GetAllCompilationTypes.
func (mr *MockCompilationMockRecorder) GetAllCompilationTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCompilationTypes", reflect.TypeOf((*MockCompilation)(nil).GetAllCompilationTypes))
}

// GetCompilation mocks base method.
func (m *MockCompilation) GetCompilation(id int) (*entity.Compilation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompilation", id)
	ret0, _ := ret[0].(*entity.Compilation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompilation indicates an expected call of GetCompilation.
func (mr *MockCompilationMockRecorder) GetCompilation(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompilation", reflect.TypeOf((*MockCompilation)(nil).GetCompilation), id)
}

// GetCompilationContent mocks base method.
func (m *MockCompilation) GetCompilationContent(id, page, limit int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompilationContent", id, page, limit)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompilationContent indicates an expected call of GetCompilationContent.
func (mr *MockCompilationMockRecorder) GetCompilationContent(id, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompilationContent", reflect.TypeOf((*MockCompilation)(nil).GetCompilationContent), id, page, limit)
}

// GetCompilationContentLength mocks base method.
func (m *MockCompilation) GetCompilationContentLength(id int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompilationContentLength", id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompilationContentLength indicates an expected call of GetCompilationContentLength.
func (mr *MockCompilationMockRecorder) GetCompilationContentLength(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompilationContentLength", reflect.TypeOf((*MockCompilation)(nil).GetCompilationContentLength), id)
}

// GetCompilationsByTypeID mocks base method.
func (m *MockCompilation) GetCompilationsByTypeID(compilationTypeID int) ([]*entity.Compilation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompilationsByTypeID", compilationTypeID)
	ret0, _ := ret[0].([]*entity.Compilation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompilationsByTypeID indicates an expected call of GetCompilationsByTypeID.
func (mr *MockCompilationMockRecorder) GetCompilationsByTypeID(compilationTypeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompilationsByTypeID", reflect.TypeOf((*MockCompilation)(nil).GetCompilationsByTypeID), compilationTypeID)
}
