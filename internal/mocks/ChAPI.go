// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	perun "github.com/hyperledger-labs/perun-node"
	mock "github.com/stretchr/testify/mock"
)

// ChAPI is an autogenerated mock type for the ChAPI type
type ChAPI struct {
	mock.Mock
}

// ChallengeDurSecs provides a mock function with given fields:
func (_m *ChAPI) ChallengeDurSecs() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// Close provides a mock function with given fields: _a0
func (_m *ChAPI) Close(_a0 context.Context) (perun.ChInfo, perun.APIError) {
	ret := _m.Called(_a0)

	var r0 perun.ChInfo
	if rf, ok := ret.Get(0).(func(context.Context) perun.ChInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(perun.ChInfo)
	}

	var r1 perun.APIError
	if rf, ok := ret.Get(1).(func(context.Context) perun.APIError); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(perun.APIError)
		}
	}

	return r0, r1
}

// Currency provides a mock function with given fields: symbol
func (_m *ChAPI) Currency(symbol string) (int, perun.Currency, bool) {
	ret := _m.Called(symbol)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(symbol)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 perun.Currency
	if rf, ok := ret.Get(1).(func(string) perun.Currency); ok {
		r1 = rf(symbol)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(perun.Currency)
		}
	}

	var r2 bool
	if rf, ok := ret.Get(2).(func(string) bool); ok {
		r2 = rf(symbol)
	} else {
		r2 = ret.Get(2).(bool)
	}

	return r0, r1, r2
}

// GetChInfo provides a mock function with given fields:
func (_m *ChAPI) GetChInfo() perun.ChInfo {
	ret := _m.Called()

	var r0 perun.ChInfo
	if rf, ok := ret.Get(0).(func() perun.ChInfo); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(perun.ChInfo)
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *ChAPI) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Parts provides a mock function with given fields:
func (_m *ChAPI) Parts() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// RespondChUpdate provides a mock function with given fields: _a0, _a1, _a2
func (_m *ChAPI) RespondChUpdate(_a0 context.Context, _a1 string, _a2 bool) (perun.ChInfo, perun.APIError) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 perun.ChInfo
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) perun.ChInfo); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(perun.ChInfo)
	}

	var r1 perun.APIError
	if rf, ok := ret.Get(1).(func(context.Context, string, bool) perun.APIError); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(perun.APIError)
		}
	}

	return r0, r1
}

// SendChUpdate provides a mock function with given fields: _a0, _a1
func (_m *ChAPI) SendChUpdate(_a0 context.Context, _a1 perun.StateUpdater) (perun.ChInfo, perun.APIError) {
	ret := _m.Called(_a0, _a1)

	var r0 perun.ChInfo
	if rf, ok := ret.Get(0).(func(context.Context, perun.StateUpdater) perun.ChInfo); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(perun.ChInfo)
	}

	var r1 perun.APIError
	if rf, ok := ret.Get(1).(func(context.Context, perun.StateUpdater) perun.APIError); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(perun.APIError)
		}
	}

	return r0, r1
}

// SubChUpdates provides a mock function with given fields: _a0
func (_m *ChAPI) SubChUpdates(_a0 perun.ChUpdateNotifier) perun.APIError {
	ret := _m.Called(_a0)

	var r0 perun.APIError
	if rf, ok := ret.Get(0).(func(perun.ChUpdateNotifier) perun.APIError); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(perun.APIError)
		}
	}

	return r0
}

// UnsubChUpdates provides a mock function with given fields:
func (_m *ChAPI) UnsubChUpdates() perun.APIError {
	ret := _m.Called()

	var r0 perun.APIError
	if rf, ok := ret.Get(0).(func() perun.APIError); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(perun.APIError)
		}
	}

	return r0
}
