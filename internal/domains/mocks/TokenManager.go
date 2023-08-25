// Code generated by mockery v2.33.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TokenManager is an autogenerated mock type for the TokenManager type
type TokenManager struct {
	mock.Mock
}

type TokenManager_Expecter struct {
	mock *mock.Mock
}

func (_m *TokenManager) EXPECT() *TokenManager_Expecter {
	return &TokenManager_Expecter{mock: &_m.Mock}
}

// GetTokens provides a mock function with given fields: ctx, guid
func (_m *TokenManager) GetTokens(ctx context.Context, guid string) (string, string, error) {
	ret := _m.Called(ctx, guid)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, string, error)); ok {
		return rf(ctx, guid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, guid)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, guid)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, guid)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TokenManager_GetTokens_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTokens'
type TokenManager_GetTokens_Call struct {
	*mock.Call
}

// GetTokens is a helper method to define mock.On call
//   - ctx context.Context
//   - guid string
func (_e *TokenManager_Expecter) GetTokens(ctx interface{}, guid interface{}) *TokenManager_GetTokens_Call {
	return &TokenManager_GetTokens_Call{Call: _e.mock.On("GetTokens", ctx, guid)}
}

func (_c *TokenManager_GetTokens_Call) Run(run func(ctx context.Context, guid string)) *TokenManager_GetTokens_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *TokenManager_GetTokens_Call) Return(_a0 string, _a1 string, _a2 error) *TokenManager_GetTokens_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *TokenManager_GetTokens_Call) RunAndReturn(run func(context.Context, string) (string, string, error)) *TokenManager_GetTokens_Call {
	_c.Call.Return(run)
	return _c
}

// RefreshTokens provides a mock function with given fields: ctx, refresh
func (_m *TokenManager) RefreshTokens(ctx context.Context, refresh string) (string, string, error) {
	ret := _m.Called(ctx, refresh)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, string, error)); ok {
		return rf(ctx, refresh)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, refresh)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, refresh)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, refresh)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TokenManager_RefreshTokens_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RefreshTokens'
type TokenManager_RefreshTokens_Call struct {
	*mock.Call
}

// RefreshTokens is a helper method to define mock.On call
//   - ctx context.Context
//   - refresh string
func (_e *TokenManager_Expecter) RefreshTokens(ctx interface{}, refresh interface{}) *TokenManager_RefreshTokens_Call {
	return &TokenManager_RefreshTokens_Call{Call: _e.mock.On("RefreshTokens", ctx, refresh)}
}

func (_c *TokenManager_RefreshTokens_Call) Run(run func(ctx context.Context, refresh string)) *TokenManager_RefreshTokens_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *TokenManager_RefreshTokens_Call) Return(_a0 string, _a1 string, _a2 error) *TokenManager_RefreshTokens_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *TokenManager_RefreshTokens_Call) RunAndReturn(run func(context.Context, string) (string, string, error)) *TokenManager_RefreshTokens_Call {
	_c.Call.Return(run)
	return _c
}

// NewTokenManager creates a new instance of TokenManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenManager {
	mock := &TokenManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
