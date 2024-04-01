// Code generated by MockGen. DO NOT EDIT.
// Source: review.go
//
// Generated by this command:
//
//	mockgen -source=review.go -destination=mocks/mock_review.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockReview is a mock of Review interface.
type MockReview struct {
	ctrl     *gomock.Controller
	recorder *MockReviewMockRecorder
}

// MockReviewMockRecorder is the mock recorder for MockReview.
type MockReviewMockRecorder struct {
	mock *MockReview
}

// NewMockReview creates a new mock instance.
func NewMockReview(ctrl *gomock.Controller) *MockReview {
	mock := &MockReview{ctrl: ctrl}
	mock.recorder = &MockReviewMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReview) EXPECT() *MockReviewMockRecorder {
	return m.recorder
}

// CreateReview mocks base method.
func (m *MockReview) CreateReview(arg0 dto.Review) (*dto.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReview", arg0)
	ret0, _ := ret[0].(*dto.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReview indicates an expected call of CreateReview.
func (mr *MockReviewMockRecorder) CreateReview(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReview", reflect.TypeOf((*MockReview)(nil).CreateReview), arg0)
}

// DeleteReview mocks base method.
func (m *MockReview) DeleteReview(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteReview", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteReview indicates an expected call of DeleteReview.
func (mr *MockReviewMockRecorder) DeleteReview(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteReview", reflect.TypeOf((*MockReview)(nil).DeleteReview), arg0)
}

// EditReview mocks base method.
func (m *MockReview) EditReview(arg0 dto.Review) (*dto.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditReview", arg0)
	ret0, _ := ret[0].(*dto.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditReview indicates an expected call of EditReview.
func (mr *MockReviewMockRecorder) EditReview(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditReview", reflect.TypeOf((*MockReview)(nil).EditReview), arg0)
}

// GetContentReviews mocks base method.
func (m *MockReview) GetContentReviews(arg0 int) (*[]dto.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentReviews", arg0)
	ret0, _ := ret[0].(*[]dto.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContentReviews indicates an expected call of GetContentReviews.
func (mr *MockReviewMockRecorder) GetContentReviews(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentReviews", reflect.TypeOf((*MockReview)(nil).GetContentReviews), arg0)
}

// GetGlobalReviews mocks base method.
func (m *MockReview) GetGlobalReviews(arg0 int) (*[]dto.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGlobalReviews", arg0)
	ret0, _ := ret[0].(*[]dto.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGlobalReviews indicates an expected call of GetGlobalReviews.
func (mr *MockReviewMockRecorder) GetGlobalReviews(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGlobalReviews", reflect.TypeOf((*MockReview)(nil).GetGlobalReviews), arg0)
}

// GetReview mocks base method.
func (m *MockReview) GetReview(arg0 int) (*dto.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReview", arg0)
	ret0, _ := ret[0].(*dto.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReview indicates an expected call of GetReview.
func (mr *MockReviewMockRecorder) GetReview(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReview", reflect.TypeOf((*MockReview)(nil).GetReview), arg0)
}

// GetUserReviews mocks base method.
func (m *MockReview) GetUserReviews(arg0, arg1 int) (*[]dto.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserReviews", arg0, arg1)
	ret0, _ := ret[0].(*[]dto.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserReviews indicates an expected call of GetUserReviews.
func (mr *MockReviewMockRecorder) GetUserReviews(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserReviews", reflect.TypeOf((*MockReview)(nil).GetUserReviews), arg0, arg1)
}
