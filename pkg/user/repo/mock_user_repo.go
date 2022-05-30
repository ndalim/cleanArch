// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/user/repo/user_repo.go

// Package repo is a generated GoMock package.
package repo

import (
	user "redditapp/pkg/user"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserRepoInterface is a mock of UserRepoInterface interface.
type MockUserRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoInterfaceMockRecorder
}

// MockUserRepoInterfaceMockRecorder is the mock recorder for MockUserRepoInterface.
type MockUserRepoInterfaceMockRecorder struct {
	mock *MockUserRepoInterface
}

// NewMockUserRepoInterface creates a new mock instance.
func NewMockUserRepoInterface(ctrl *gomock.Controller) *MockUserRepoInterface {
	mock := &MockUserRepoInterface{ctrl: ctrl}
	mock.recorder = &MockUserRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepoInterface) EXPECT() *MockUserRepoInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepoInterface) Create(arg0 string, arg1 []byte) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserRepoInterfaceMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepoInterface)(nil).Create), arg0, arg1)
}

// GetAll mocks base method.
func (m *MockUserRepoInterface) GetAll() ([]*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockUserRepoInterfaceMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUserRepoInterface)(nil).GetAll))
}

// GetUser mocks base method.
func (m *MockUserRepoInterface) GetUser(name string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", name)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserRepoInterfaceMockRecorder) GetUser(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserRepoInterface)(nil).GetUser), name)
}