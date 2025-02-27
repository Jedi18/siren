// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	"github.com/odpf/siren/domain"
	mock "github.com/stretchr/testify/mock"
)

// WorkspaceService is an autogenerated mock type for the WorkspaceService type
type WorkspaceService struct {
	mock.Mock
}

// GetChannels provides a mock function with given fields: _a0
func (_m *WorkspaceService) GetChannels(_a0 string) ([]domain.Channel, error) {
	ret := _m.Called(_a0)

	var r0 []domain.Channel
	if rf, ok := ret.Get(0).(func(string) []domain.Channel); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
