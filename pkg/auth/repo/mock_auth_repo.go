// Code generated by MockGen. DO NOT EDIT.
// Source: auth_repo.go

// Package mock_repo is a generated GoMock package.
package repo

import (
	user "redditapp/pkg/user"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthRepoInterface is a mock of AuthRepoInterface interface.
type MockAuthRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepoInterfaceMockRecorder
}

// MockAuthRepoInterfaceMockRecorder is the mock recorder for MockAuthRepoInterface.
type MockAuthRepoInterfaceMockRecorder struct {
	mock *MockAuthRepoInterface
}

// NewMockAuthRepoInterface creates a new mock instance.
func NewMockAuthRepoInterface(ctrl *gomock.Controller) *MockAuthRepoInterface {
	mock := &MockAuthRepoInterface{ctrl: ctrl}
	mock.recorder = &MockAuthRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepoInterface) EXPECT() *MockAuthRepoInterfaceMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockAuthRepoInterface) Check(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockAuthRepoInterfaceMockRecorder) Check(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockAuthRepoInterface)(nil).Check), arg0)
}

// Session mocks base method.
func (m *MockAuthRepoInterface) Session(user *user.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Session", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Session indicates an expected call of Session.
func (mr *MockAuthRepoInterfaceMockRecorder) Session(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Session", reflect.TypeOf((*MockAuthRepoInterface)(nil).Session), user)
}
