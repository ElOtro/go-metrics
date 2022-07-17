// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Repo is an autogenerated mock type for the Repo type
type Repo struct {
	mock.Mock
}

// Get provides a mock function with given fields: t, n
func (_m *Repo) Get(t string, n string) (string, error) {
	ret := _m.Called(t, n)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(t, n)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(t, n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *Repo) GetAll() (map[string]float64, map[string]int64) {
	ret := _m.Called()

	var r0 map[string]float64
	if rf, ok := ret.Get(0).(func() map[string]float64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]float64)
		}
	}

	var r1 map[string]int64
	if rf, ok := ret.Get(1).(func() map[string]int64); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]int64)
		}
	}

	return r0, r1
}

// Set provides a mock function with given fields: t, n, v
func (_m *Repo) Set(t string, n string, v string) error {
	ret := _m.Called(t, n, v)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(t, n, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepo creates a new instance of Repo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepo(t mockConstructorTestingTNewRepo) *Repo {
	mock := &Repo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
