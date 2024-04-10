package limitermocks

import (
	"context"
	"reflect"

	"go.uber.org/mock/gomock"
)

type MockLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockLimiterMockRecorder
}

type MockLimiterMockRecorder struct {
	mock *MockLimiter
}

func NewMockLimiter(ctrl *gomock.Controller) *MockLimiter {
	mock := &MockLimiter{ctrl: ctrl}
	mock.recorder = &MockLimiterMockRecorder{mock}

	return mock
}

func (m *MockLimiter) EXPECT() *MockLimiterMockRecorder {
	return m.recorder
}

func (m *MockLimiter) Limit(ctx context.Context, key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Limit", ctx, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)

	return ret0, ret1
}

func (mr *MockLimiterMockRecorder) Limit(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Limit", reflect.TypeOf((*MockLimiter)(nil).Limit), ctx, key)
}
