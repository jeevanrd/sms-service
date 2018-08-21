// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import smsservice "github.com/jeevanrd/sms-service/smsservice"

// SmsService is an autogenerated mock type for the SmsService type
type SmsService struct {
	mock.Mock
}

// InboundSms provides a mock function with given fields: ctx, req
func (_m *SmsService) InboundSms(ctx context.Context, req smsservice.SmsRequest) (interface{}, error) {
	ret := _m.Called(ctx, req)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, smsservice.SmsRequest) interface{}); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, smsservice.SmsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OutboundSms provides a mock function with given fields: ctx, req
func (_m *SmsService) OutboundSms(ctx context.Context, req smsservice.SmsRequest) (interface{}, error) {
	ret := _m.Called(ctx, req)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, smsservice.SmsRequest) interface{}); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, smsservice.SmsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
