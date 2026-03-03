// Manually maintained mock for TaskEnqueuer interface (testing only)

package schemas

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTaskEnqueuer is a mock of TaskEnqueuer interface.
type MockTaskEnqueuer struct {
	ctrl     *gomock.Controller
	recorder *MockTaskEnqueuerMockRecorder
}

// MockTaskEnqueuerMockRecorder is the mock recorder for MockTaskEnqueuer.
type MockTaskEnqueuerMockRecorder struct {
	mock *MockTaskEnqueuer
}

// NewMockTaskEnqueuer creates a new mock instance.
func NewMockTaskEnqueuer(ctrl *gomock.Controller) *MockTaskEnqueuer {
	mock := &MockTaskEnqueuer{ctrl: ctrl}
	mock.recorder = &MockTaskEnqueuerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskEnqueuer) EXPECT() *MockTaskEnqueuerMockRecorder {
	return m.recorder
}

// EnqueueTransactionWithPriority mocks base method.
func (m *MockTaskEnqueuer) EnqueueTransactionWithPriority(ctx context.Context, event map[string]any, priority string) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnqueueTransactionWithPriority", ctx, event, priority)
	ret0 := ret[0]
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnqueueTransactionWithPriority indicates an expected call of EnqueueTransactionWithPriority.
func (mr *MockTaskEnqueuerMockRecorder) EnqueueTransactionWithPriority(ctx, event, priority any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueTransactionWithPriority", reflect.TypeOf((*MockTaskEnqueuer)(nil).EnqueueTransactionWithPriority), ctx, event, priority)
}
