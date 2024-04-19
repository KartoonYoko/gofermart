// Code generated by MockGen. DO NOT EDIT.
// Source: gofermart/internal/usecase/order_accrual (interfaces: OrderAccrualRepository,OrderAccrualAPI)
//
// Generated by this command:
//
//	mockgen --destination=internal/mocks/order_accrual/mock.go --package=order_accrual gofermart/internal/usecase/order_accrual OrderAccrualRepository,OrderAccrualAPI
//

// Package order_accrual is a generated GoMock package.
package order_accrual

import (
	context "context"
	auth "gofermart/internal/model/auth"
	order "gofermart/internal/model/order"
	orderaccrual "gofermart/internal/model/order_accrual"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockOrderAccrualRepository is a mock of OrderAccrualRepository interface.
type MockOrderAccrualRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderAccrualRepositoryMockRecorder
}

// MockOrderAccrualRepositoryMockRecorder is the mock recorder for MockOrderAccrualRepository.
type MockOrderAccrualRepositoryMockRecorder struct {
	mock *MockOrderAccrualRepository
}

// NewMockOrderAccrualRepository creates a new mock instance.
func NewMockOrderAccrualRepository(ctrl *gomock.Controller) *MockOrderAccrualRepository {
	mock := &MockOrderAccrualRepository{ctrl: ctrl}
	mock.recorder = &MockOrderAccrualRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderAccrualRepository) EXPECT() *MockOrderAccrualRepositoryMockRecorder {
	return m.recorder
}

// AccrualToOrderAndUser mocks base method.
func (m *MockOrderAccrualRepository) AccrualToOrderAndUser(arg0 context.Context, arg1 int64, arg2 auth.UserID, arg3 int, arg4 order.OrderStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccrualToOrderAndUser", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// AccrualToOrderAndUser indicates an expected call of AccrualToOrderAndUser.
func (mr *MockOrderAccrualRepositoryMockRecorder) AccrualToOrderAndUser(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccrualToOrderAndUser", reflect.TypeOf((*MockOrderAccrualRepository)(nil).AccrualToOrderAndUser), arg0, arg1, arg2, arg3, arg4)
}

// GetUnhandledOrders mocks base method.
func (m *MockOrderAccrualRepository) GetUnhandledOrders(arg0 context.Context) ([]orderaccrual.GetOrderModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnhandledOrders", arg0)
	ret0, _ := ret[0].([]orderaccrual.GetOrderModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnhandledOrders indicates an expected call of GetUnhandledOrders.
func (mr *MockOrderAccrualRepositoryMockRecorder) GetUnhandledOrders(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnhandledOrders", reflect.TypeOf((*MockOrderAccrualRepository)(nil).GetUnhandledOrders), arg0)
}

// MockOrderAccrualAPI is a mock of OrderAccrualAPI interface.
type MockOrderAccrualAPI struct {
	ctrl     *gomock.Controller
	recorder *MockOrderAccrualAPIMockRecorder
}

// MockOrderAccrualAPIMockRecorder is the mock recorder for MockOrderAccrualAPI.
type MockOrderAccrualAPIMockRecorder struct {
	mock *MockOrderAccrualAPI
}

// NewMockOrderAccrualAPI creates a new mock instance.
func NewMockOrderAccrualAPI(ctrl *gomock.Controller) *MockOrderAccrualAPI {
	mock := &MockOrderAccrualAPI{ctrl: ctrl}
	mock.recorder = &MockOrderAccrualAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderAccrualAPI) EXPECT() *MockOrderAccrualAPIMockRecorder {
	return m.recorder
}

// GetOrderAccrual mocks base method.
func (m *MockOrderAccrualAPI) GetOrderAccrual(arg0 context.Context, arg1 int64) (*orderaccrual.GetOrderAccrualFromRemoteModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderAccrual", arg0, arg1)
	ret0, _ := ret[0].(*orderaccrual.GetOrderAccrualFromRemoteModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderAccrual indicates an expected call of GetOrderAccrual.
func (mr *MockOrderAccrualAPIMockRecorder) GetOrderAccrual(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderAccrual", reflect.TypeOf((*MockOrderAccrualAPI)(nil).GetOrderAccrual), arg0, arg1)
}
