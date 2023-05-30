// Code generated by mockery v2.28.1. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"

// MockResult is an autogenerated mock type for the Result type
type MockResult struct {
	mock.Mock
}

type MockResult_Expecter struct {
	mock *mock.Mock
}

func (_m *MockResult) EXPECT() *MockResult_Expecter {
	return &MockResult_Expecter{mock: &_m.Mock}
}

// GetVal provides a mock function with given fields:
func (_m *MockResult) GetVal() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockResult_GetVal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVal'
type MockResult_GetVal_Call struct {
	*mock.Call
}

// GetVal is a helper method to define mock.On call
func (_e *MockResult_Expecter) GetVal() *MockResult_GetVal_Call {
	return &MockResult_GetVal_Call{Call: _e.mock.On("GetVal")}
}

func (_c *MockResult_GetVal_Call) Run(run func()) *MockResult_GetVal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResult_GetVal_Call) Return(_a0 int64) *MockResult_GetVal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResult_GetVal_Call) RunAndReturn(run func() int64) *MockResult_GetVal_Call {
	_c.Call.Return(run)
	return _c
}

// ToString provides a mock function with given fields:
func (_m *MockResult) ToString() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockResult_ToString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ToString'
type MockResult_ToString_Call struct {
	*mock.Call
}

// ToString is a helper method to define mock.On call
func (_e *MockResult_Expecter) ToString() *MockResult_ToString_Call {
	return &MockResult_ToString_Call{Call: _e.mock.On("ToString")}
}

func (_c *MockResult_ToString_Call) Run(run func()) *MockResult_ToString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResult_ToString_Call) Return(_a0 string) *MockResult_ToString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResult_ToString_Call) RunAndReturn(run func() string) *MockResult_ToString_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockResult interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockResult creates a new instance of MockResult. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockResult(t mockConstructorTestingTNewMockResult) *MockResult {
	mock := &MockResult{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
