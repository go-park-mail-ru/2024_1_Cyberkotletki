// Code generated by MockGen. DO NOT EDIT.
// Source: compilation.go
//
// Generated by this command:
//
//	mockgen -source=compilation.go -destination=mocks/mock_compilation.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
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

// GetCompilationContent mocks base method.
func (m *MockCompilation) GetCompilationContent(compID, page, limit int) ([]*dto.PreviewContentCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompilationContent", compID, page, limit)
	ret0, _ := ret[0].([]*dto.PreviewContentCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompilationContent indicates an expected call of GetCompilationContent.
func (mr *MockCompilationMockRecorder) GetCompilationContent(compID, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompilationContent", reflect.TypeOf((*MockCompilation)(nil).GetCompilationContent), compID, page, limit)
}

// GetCompilationTypes mocks base method.
func (m *MockCompilation) GetCompilationTypes() (*dto.CompilationTypeResponseList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompilationTypes")
	ret0, _ := ret[0].(*dto.CompilationTypeResponseList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompilationTypes indicates an expected call of GetCompilationTypes.
func (mr *MockCompilationMockRecorder) GetCompilationTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompilationTypes", reflect.TypeOf((*MockCompilation)(nil).GetCompilationTypes))
}

// GetCompilationsByCompilationType mocks base method.
func (m *MockCompilation) GetCompilationsByCompilationType(compTypeID int) (*dto.CompilationResponseList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompilationsByCompilationType", compTypeID)
	ret0, _ := ret[0].(*dto.CompilationResponseList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompilationsByCompilationType indicates an expected call of GetCompilationsByCompilationType.
func (mr *MockCompilationMockRecorder) GetCompilationsByCompilationType(compTypeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompilationsByCompilationType", reflect.TypeOf((*MockCompilation)(nil).GetCompilationsByCompilationType), compTypeID)
}
