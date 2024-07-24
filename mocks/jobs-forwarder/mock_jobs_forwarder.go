// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rudderlabs/rudder-server/schema-forwarder (interfaces: Forwarder)
//
// Generated by this command:
//
//	mockgen -destination=../mocks/jobs-forwarder/mock_jobs_forwarder.go -package=mock_jobs_forwarder github.com/rudderlabs/rudder-server/schema-forwarder Forwarder
//

// Package mock_jobs_forwarder is a generated GoMock package.
package mock_jobs_forwarder

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockForwarder is a mock of Forwarder interface.
type MockForwarder struct {
	ctrl     *gomock.Controller
	recorder *MockForwarderMockRecorder
}

// MockForwarderMockRecorder is the mock recorder for MockForwarder.
type MockForwarderMockRecorder struct {
	mock *MockForwarder
}

// NewMockForwarder creates a new mock instance.
func NewMockForwarder(ctrl *gomock.Controller) *MockForwarder {
	mock := &MockForwarder{ctrl: ctrl}
	mock.recorder = &MockForwarderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockForwarder) EXPECT() *MockForwarderMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockForwarder) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockForwarderMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockForwarder)(nil).Start))
}

// Stop mocks base method.
func (m *MockForwarder) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockForwarderMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockForwarder)(nil).Stop))
}
