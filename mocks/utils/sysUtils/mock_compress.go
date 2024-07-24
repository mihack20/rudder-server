// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rudderlabs/rudder-server/utils/sysUtils (interfaces: GZipI)
//
// Generated by this command:
//
//	mockgen -destination=../../mocks/utils/sysUtils/mock_compress.go -package mock_sysUtils github.com/rudderlabs/rudder-server/utils/sysUtils GZipI
//

// Package mock_sysUtils is a generated GoMock package.
package mock_sysUtils

import (
	gzip "compress/gzip"
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockGZipI is a mock of GZipI interface.
type MockGZipI struct {
	ctrl     *gomock.Controller
	recorder *MockGZipIMockRecorder
}

// MockGZipIMockRecorder is the mock recorder for MockGZipI.
type MockGZipIMockRecorder struct {
	mock *MockGZipI
}

// NewMockGZipI creates a new mock instance.
func NewMockGZipI(ctrl *gomock.Controller) *MockGZipI {
	mock := &MockGZipI{ctrl: ctrl}
	mock.recorder = &MockGZipIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGZipI) EXPECT() *MockGZipIMockRecorder {
	return m.recorder
}

// NewReader mocks base method.
func (m *MockGZipI) NewReader(arg0 io.Reader) (*gzip.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewReader", arg0)
	ret0, _ := ret[0].(*gzip.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewReader indicates an expected call of NewReader.
func (mr *MockGZipIMockRecorder) NewReader(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewReader", reflect.TypeOf((*MockGZipI)(nil).NewReader), arg0)
}

// NewWriter mocks base method.
func (m *MockGZipI) NewWriter(arg0 io.Writer) *gzip.Writer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWriter", arg0)
	ret0, _ := ret[0].(*gzip.Writer)
	return ret0
}

// NewWriter indicates an expected call of NewWriter.
func (mr *MockGZipIMockRecorder) NewWriter(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWriter", reflect.TypeOf((*MockGZipI)(nil).NewWriter), arg0)
}
