// Code generated by mockery v2.9.4. DO NOT EDIT.

package transaction

import mock "github.com/stretchr/testify/mock"

// errorer is an autogenerated mock type for the errorer type
type errorer struct {
	mock.Mock
}

// error provides a mock function with given fields:
func (_m *errorer) error() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
