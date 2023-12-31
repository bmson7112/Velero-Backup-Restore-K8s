// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	v1 "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned/typed/velero/v1"
)

// SchedulesGetter is an autogenerated mock type for the SchedulesGetter type
type SchedulesGetter struct {
	mock.Mock
}

// Schedules provides a mock function with given fields: namespace
func (_m *SchedulesGetter) Schedules(namespace string) v1.ScheduleInterface {
	ret := _m.Called(namespace)

	var r0 v1.ScheduleInterface
	if rf, ok := ret.Get(0).(func(string) v1.ScheduleInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ScheduleInterface)
		}
	}

	return r0
}

// NewSchedulesGetter creates a new instance of SchedulesGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSchedulesGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *SchedulesGetter {
	mock := &SchedulesGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
