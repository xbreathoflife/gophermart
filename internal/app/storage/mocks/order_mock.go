// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/xbreathoflife/gophermart/internal/app/entities"
)

// MockOrderStorage is a mock of OrderStorage interface.
type MockOrderStorage struct {
	ctrl     *gomock.Controller
	recorder *MockOrderStorageMockRecorder
}

// MockOrderStorageMockRecorder is the mock recorder for MockOrderStorage.
type MockOrderStorageMockRecorder struct {
	mock *MockOrderStorage
}

// NewMockOrderStorage creates a new mock instance.
func NewMockOrderStorage(ctrl *gomock.Controller) *MockOrderStorage {
	mock := &MockOrderStorage{ctrl: ctrl}
	mock.recorder = &MockOrderStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderStorage) EXPECT() *MockOrderStorageMockRecorder {
	return m.recorder
}

// GetOrderIfExists mocks base method.
func (m *MockOrderStorage) GetOrderIfExists(ctx context.Context, orderNum string) (*entities.OrderModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderIfExists", ctx, orderNum)
	ret0, _ := ret[0].(*entities.OrderModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderIfExists indicates an expected call of GetOrderIfExists.
func (mr *MockOrderStorageMockRecorder) GetOrderIfExists(ctx, orderNum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderIfExists", reflect.TypeOf((*MockOrderStorage)(nil).GetOrderIfExists), ctx, orderNum)
}

// GetOrdersForUser mocks base method.
func (m *MockOrderStorage) GetOrdersForUser(ctx context.Context, login string) ([]entities.OrderModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersForUser", ctx, login)
	ret0, _ := ret[0].([]entities.OrderModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersForUser indicates an expected call of GetOrdersForUser.
func (mr *MockOrderStorageMockRecorder) GetOrdersForUser(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersForUser", reflect.TypeOf((*MockOrderStorage)(nil).GetOrdersForUser), ctx, login)
}

// InsertNewOrder mocks base method.
func (m *MockOrderStorage) InsertNewOrder(ctx context.Context, order entities.OrderModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewOrder indicates an expected call of InsertNewOrder.
func (mr *MockOrderStorageMockRecorder) InsertNewOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewOrder", reflect.TypeOf((*MockOrderStorage)(nil).InsertNewOrder), ctx, order)
}

// UpdateOrderStatus mocks base method.
func (m *MockOrderStorage) UpdateOrderStatus(ctx context.Context, orderNum, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", ctx, orderNum, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockOrderStorageMockRecorder) UpdateOrderStatus(ctx, orderNum, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockOrderStorage)(nil).UpdateOrderStatus), ctx, orderNum, status)
}

// UpdateOrderStatusAndAccrual mocks base method.
func (m *MockOrderStorage) UpdateOrderStatusAndAccrual(ctx context.Context, orderNum, status string, accrual float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatusAndAccrual", ctx, orderNum, status, accrual)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatusAndAccrual indicates an expected call of UpdateOrderStatusAndAccrual.
func (mr *MockOrderStorageMockRecorder) UpdateOrderStatusAndAccrual(ctx, orderNum, status, accrual interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatusAndAccrual", reflect.TypeOf((*MockOrderStorage)(nil).UpdateOrderStatusAndAccrual), ctx, orderNum, status, accrual)
}
