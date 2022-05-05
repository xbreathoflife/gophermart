// Code generated by MockGen. DO NOT EDIT.
// Source: balance.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/xbreathoflife/gophermart/internal/app/entities"
)

// MockBalanceStorage is a mock of BalanceStorage interface.
type MockBalanceStorage struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceStorageMockRecorder
}

// MockBalanceStorageMockRecorder is the mock recorder for MockBalanceStorage.
type MockBalanceStorageMockRecorder struct {
	mock *MockBalanceStorage
}

// NewMockBalanceStorage creates a new mock instance.
func NewMockBalanceStorage(ctrl *gomock.Controller) *MockBalanceStorage {
	mock := &MockBalanceStorage{ctrl: ctrl}
	mock.recorder = &MockBalanceStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalanceStorage) EXPECT() *MockBalanceStorageMockRecorder {
	return m.recorder
}

// GetBalance mocks base method.
func (m *MockBalanceStorage) GetBalance(ctx context.Context, login string) (*entities.BalanceModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, login)
	ret0, _ := ret[0].(*entities.BalanceModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockBalanceStorageMockRecorder) GetBalance(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockBalanceStorage)(nil).GetBalance), ctx, login)
}

// GetBalanceWithdrawalsForUser mocks base method.
func (m *MockBalanceStorage) GetBalanceWithdrawalsForUser(ctx context.Context, login string) ([]entities.BalanceWithdrawalsModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalanceWithdrawalsForUser", ctx, login)
	ret0, _ := ret[0].([]entities.BalanceWithdrawalsModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalanceWithdrawalsForUser indicates an expected call of GetBalanceWithdrawalsForUser.
func (mr *MockBalanceStorageMockRecorder) GetBalanceWithdrawalsForUser(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalanceWithdrawalsForUser", reflect.TypeOf((*MockBalanceStorage)(nil).GetBalanceWithdrawalsForUser), ctx, login)
}

// InsertNewBalance mocks base method.
func (m *MockBalanceStorage) InsertNewBalance(ctx context.Context, balance entities.BalanceModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewBalance", ctx, balance)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewBalance indicates an expected call of InsertNewBalance.
func (mr *MockBalanceStorageMockRecorder) InsertNewBalance(ctx, balance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewBalance", reflect.TypeOf((*MockBalanceStorage)(nil).InsertNewBalance), ctx, balance)
}

// InsertNewBalanceWithdrawals mocks base method.
func (m *MockBalanceStorage) InsertNewBalanceWithdrawals(ctx context.Context, balanceWithdrawals entities.BalanceWithdrawalsModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewBalanceWithdrawals", ctx, balanceWithdrawals)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewBalanceWithdrawals indicates an expected call of InsertNewBalanceWithdrawals.
func (mr *MockBalanceStorageMockRecorder) InsertNewBalanceWithdrawals(ctx, balanceWithdrawals interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewBalanceWithdrawals", reflect.TypeOf((*MockBalanceStorage)(nil).InsertNewBalanceWithdrawals), ctx, balanceWithdrawals)
}

// UpdateBalance mocks base method.
func (m *MockBalanceStorage) UpdateBalance(ctx context.Context, balance entities.BalanceModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", ctx, balance)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockBalanceStorageMockRecorder) UpdateBalance(ctx, balance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockBalanceStorage)(nil).UpdateBalance), ctx, balance)
}