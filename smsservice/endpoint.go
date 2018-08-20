package smsservice

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type SmsRequest struct {
	From string
	To string
	Text string
}

func makeInboundSmsEndpoint(service SmsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SmsRequest)
		i, err := service.InboundSms(ctx,req)
		return i, err
	}
}

func makeOutboundSmsEndpoint(service SmsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SmsRequest)
		i, err := service.OutboundSms(ctx, req)
		return i, err
	}
}
