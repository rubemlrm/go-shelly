// Code generated by mockery v2.31.0. DO NOT EDIT.

package client_mocks

import (
	http "net/http"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	mock "github.com/stretchr/testify/mock"
)

// MockClientProxy is an autogenerated mock type for the ClientProxy type
type MockClientProxy struct {
	mock.Mock
}

type MockClientProxy_Expecter struct {
	mock *mock.Mock
}

func (_m *MockClientProxy) EXPECT() *MockClientProxy_Expecter {
	return &MockClientProxy_Expecter{mock: &_m.Mock}
}

// Do provides a mock function with given fields: req
func (_m *MockClientProxy) Do(req *retryablehttp.Request) (*http.Response, error) {
	ret := _m.Called(req)

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(*retryablehttp.Request) (*http.Response, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*retryablehttp.Request) *http.Response); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(*retryablehttp.Request) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClientProxy_Do_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Do'
type MockClientProxy_Do_Call struct {
	*mock.Call
}

// Do is a helper method to define mock.On call
//   - req *retryablehttp.Request
func (_e *MockClientProxy_Expecter) Do(req interface{}) *MockClientProxy_Do_Call {
	return &MockClientProxy_Do_Call{Call: _e.mock.On("Do", req)}
}

func (_c *MockClientProxy_Do_Call) Run(run func(req *retryablehttp.Request)) *MockClientProxy_Do_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*retryablehttp.Request))
	})
	return _c
}

func (_c *MockClientProxy_Do_Call) Return(_a0 *http.Response, _a1 error) *MockClientProxy_Do_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClientProxy_Do_Call) RunAndReturn(run func(*retryablehttp.Request) (*http.Response, error)) *MockClientProxy_Do_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockClientProxy creates a new instance of MockClientProxy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClientProxy(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClientProxy {
	mock := &MockClientProxy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}