package smsservice

import (
	"github.com/gorilla/mux"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"io/ioutil"
	"context"
	"encoding/json"
	"errors"
	"github.com/jeevanrd/sms-service/statusErrors"
)

func MakeHandler(is SmsService, logger kitlog.Logger, r *mux.Router, accessControl func(http.Handler) http.Handler, contentTypeHandle func(http.Handler) http.Handler) *mux.Router {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(EncodeError),
	}

	inBoundSmsHandler := kithttp.NewServer(
		makeInboundSmsEndpoint(is),
		decodeInboundSms,
		encodeResponse,
		opts...,
	)

	outBoundSmsHandler := kithttp.NewServer(
		makeOutboundSmsEndpoint(is),
		decodeOutboundSms,
		encodeResponse,
		opts...,
	)

	r.Handle("/inbound/sms", contentTypeHandle(accessControl(inBoundSmsHandler)))
	r.Handle("/outbound/sms", contentTypeHandle(accessControl(outBoundSmsHandler)))

	return r
}

func decodeInboundSms(c context.Context, r *http.Request) (interface{}, error) {
	if(r.Method != "POST") {
		return SmsRequest{}, statusErrors.New(errors.New("Method Not Allowed"), 405);
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return validateRequest(data)
}

func decodeOutboundSms(c context.Context, r *http.Request) (interface{}, error) {
	if(r.Method != "POST") {
		return SmsRequest{}, statusErrors.New(errors.New("Method Not Allowed"), 405);
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return validateRequest(data)
}

//validate from,to,text params
func validateRequest(body []byte) (SmsRequest,error) {
	var payload map[string]string
	err := json.Unmarshal(body, &payload)

	if(err != nil) {
		return SmsRequest{}, statusErrors.New(errors.New("invalid payload"), 400)
	}

	from, exists := payload["from"]
	if(!exists) {
		return SmsRequest{}, statusErrors.New(errors.New("from is missing"), 400)
	}

	length := len(from)
	if (!(length >= 6 && length <= 16)) {
		return SmsRequest{}, statusErrors.New(errors.New("from is invalid"), 400)
	}

	to, exists := payload["to"]
	if(!exists) {
		return SmsRequest{}, statusErrors.New(errors.New("to is missing"), 400)
	}

	length = len(to)
	if (!(length >= 6 && length <= 16)) {
		return SmsRequest{}, statusErrors.New(errors.New("to is invalid"), 400)
	}

	text, exists := payload["text"]
	if(!exists) {
		return SmsRequest{}, statusErrors.New(errors.New("text is missing"), 400)
	}

	length = len(text)
	if (!(length >= 1 && length <= 120)) {
		return SmsRequest{}, statusErrors.New(errors.New("text is invalid"), 400)
	}

	return SmsRequest{From:from,To:to,Text:text}, nil
}

//encodeResponse encodes response
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	ba, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(ba)
	w.WriteHeader(200)
	return err
}

//EncodeError encode errors from business-logic
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch e := err.(type) {
		case statusErrors.StatusError:
		switch e.StatusCode {
			case 400:
				w.WriteHeader(e.StatusCode)
				ba,_ := json.Marshal(Response{Error:err.Error(), Message:""})
				w.Write(ba)
				return
			case 404:
				w.WriteHeader(e.StatusCode)
				ba,_ := json.Marshal(Response{Error:err.Error(), Message:""})
				w.Write(ba)
				return
			case 405:
				w.WriteHeader(e.StatusCode)
				ba,_ := json.Marshal(Response{Error:err.Error(), Message:""})
				w.Write(ba)
				return
			case 429:
				w.WriteHeader(e.StatusCode)
				ba,_ := json.Marshal(Response{Error:err.Error(), Message:""})
				w.Write(ba)
				return
			default:
				w.WriteHeader(500)
				ba,_ := json.Marshal(Response{Error:"unknown failure", Message:""})
				w.Write(ba)
				return
			}
	default:
		w.WriteHeader(500)
		ba,_ := json.Marshal(Response{Error:"unknown failure", Message:""})
		w.Write(ba)
		return
	}
}