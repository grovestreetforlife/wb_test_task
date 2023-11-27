// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package mock_v1api is a generated GoMock package.
package mock_v1api

import (
	context "context"
	reflect "reflect"
	model "wb_test_task/libs/model"

	gomock "github.com/golang/mock/gomock"
)

// MockorderService is a mock of orderService interface.
type MockorderService struct {
	ctrl     *gomock.Controller
	recorder *MockorderServiceMockRecorder
}

// MockorderServiceMockRecorder is the mock recorder for MockorderService.
type MockorderServiceMockRecorder struct {
	mock *MockorderService
}

// NewMockorderService creates a new mock instance.
func NewMockorderService(ctrl *gomock.Controller) *MockorderService {
	mock := &MockorderService{ctrl: ctrl}
	mock.recorder = &MockorderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockorderService) EXPECT() *MockorderServiceMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockorderService) GetByID(ctx context.Context, id string) (*model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockorderServiceMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockorderService)(nil).GetByID), ctx, id)
}