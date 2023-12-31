// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"
	model "wb_test_task/libs/model"

	gomock "github.com/golang/mock/gomock"
)

// MockorderStorage is a mock of orderStorage interface.
type MockorderStorage struct {
	ctrl     *gomock.Controller
	recorder *MockorderStorageMockRecorder
}

// MockorderStorageMockRecorder is the mock recorder for MockorderStorage.
type MockorderStorageMockRecorder struct {
	mock *MockorderStorage
}

// NewMockorderStorage creates a new mock instance.
func NewMockorderStorage(ctrl *gomock.Controller) *MockorderStorage {
	mock := &MockorderStorage{ctrl: ctrl}
	mock.recorder = &MockorderStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockorderStorage) EXPECT() *MockorderStorageMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockorderStorage) GetByID(ctx context.Context, id string) (*model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockorderStorageMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockorderStorage)(nil).GetByID), ctx, id)
}

// MockorderCache is a mock of orderCache interface.
type MockorderCache struct {
	ctrl     *gomock.Controller
	recorder *MockorderCacheMockRecorder
}

// MockorderCacheMockRecorder is the mock recorder for MockorderCache.
type MockorderCacheMockRecorder struct {
	mock *MockorderCache
}

// NewMockorderCache creates a new mock instance.
func NewMockorderCache(ctrl *gomock.Controller) *MockorderCache {
	mock := &MockorderCache{ctrl: ctrl}
	mock.recorder = &MockorderCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockorderCache) EXPECT() *MockorderCacheMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockorderCache) GetByID(ctx context.Context, id string) (*model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockorderCacheMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockorderCache)(nil).GetByID), ctx, id)
}

// Set mocks base method.
func (m *MockorderCache) Set(ctx context.Context, key string, order *model.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockorderCacheMockRecorder) Set(ctx, key, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockorderCache)(nil).Set), ctx, key, order)
}
