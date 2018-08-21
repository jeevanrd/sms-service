package smsservice_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/jeevanrd/sms-service/mockData"
	"github.com/jeevanrd/sms-service/smsservice/mocks"
	"github.com/jeevanrd/sms-service/smsservice"
	"bytes"
	"encoding/json"
)

func preparePostRequest(url string, body io.Reader) *http.Request {
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil
	}
	return addUserInfo(request)
}

func addUserInfo(request *http.Request) *http.Request {
	ctx := context.WithValue(request.Context(), "user-information", mockData.Account1)
	ctx = context.WithValue(ctx, "request-method", request.Method)
	request = request.WithContext(ctx)
	return request
}

func toJsonStr(obj interface{}) []byte {
	ret, err := json.Marshal(obj)
	if err != nil {
		panic("should not happen, unable to create json")
	}
	return ret
}


type serializedLogger struct {
	mtx sync.Mutex
	log.Logger
}

var _ = Describe("Transport Test", func() {

	var recorder *httptest.ResponseRecorder
	var logger log.Logger

	logger = log.NewJSONLogger(os.Stdout)
	logger = &serializedLogger{Logger: logger}
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	is := &mocks.SmsService{}
	var sampleSmsR1 = smsservice.SmsRequest{From:mockData.Phone1.Number,To:mockData.Phone2.Number,Text:"test1"}
	var sampleSmsR2 = smsservice.SmsRequest{From:mockData.Phone2.Number,To:mockData.Phone1.Number,Text:"test2"}

	is.On("InboundSms", mock.MatchedBy(func(interface{}) bool { return true }), sampleSmsR1).Return(smsservice.Response{"inbound sms ok", ""}, nil)
	is.On("OutboundSms", mock.MatchedBy(func(interface{}) bool { return true }), sampleSmsR2).Return(smsservice.Response{"outbound sms ok", ""}, nil)

	httpLogger := log.With(logger, "component", "http")
	var ac = func(h http.Handler) http.Handler {
		return h
	}

	handler := smsservice.MakeHandler(is,httpLogger, mux.NewRouter(), ac, ac)

	BeforeEach(func() {
		// Set up a new server before each test and record HTTP responses.
		recorder = httptest.NewRecorder()

	})

	Describe("POST /inbound/sms", func() {
		Context("return 400 incase of missing params", func() {
			It("should get from param error", func() {
				url := fmt.Sprintf("/inbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("from is missing"))
			})

			It("should get from param length error", func() {
				url := fmt.Sprintf("/inbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"123"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("from is invalid"))
			})

			It("should get to param error", func() {
				url := fmt.Sprintf("/inbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("to is missing"))
			})

			It("should get to param length error", func() {
				url := fmt.Sprintf("/inbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567", "to": "123"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("to is invalid"))
			})

			It("should get text param error", func() {
				url := fmt.Sprintf("/inbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567", "to": "12345666"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("text is missing"))
			})

			It("should get text param length error", func() {
				url := fmt.Sprintf("/inbound/sms")
				text := "erererererererererererererererererererhn vihnhihklneriokrnkenei"
				text = text + text
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567", "to": "12345666", "text":text})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("text is invalid"))
			})

			It("should get success", func() {
				url := fmt.Sprintf("/inbound/sms")
				text := "test1"
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":mockData.Phone1.Number, "to": mockData.Phone2.Number, "text":text})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).Should(ContainSubstring("inbound sms ok"))
			})

		})
	})

	Describe("POST /outbound/sms", func() {
		Context("return 400 incase of missing params", func() {
			It("should get from param error", func() {
				url := fmt.Sprintf("/outbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("from is missing"))
			})

			It("should get from param length error", func() {
				url := fmt.Sprintf("/outbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"123"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("from is invalid"))
			})

			It("should get to param error", func() {
				url := fmt.Sprintf("/outbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("to is missing"))
			})

			It("should get to param length error", func() {
				url := fmt.Sprintf("/outbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567", "to": "123"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("to is invalid"))
			})

			It("should get text param error", func() {
				url := fmt.Sprintf("/outbound/sms")
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567", "to": "12345666"})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("text is missing"))
			})

			It("should get text param length error", func() {
				url := fmt.Sprintf("/outbound/sms")
				text := "erererererererererererererererererererhn vihnhihklneriokrnkenei"
				text = text + text
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":"1234567", "to": "12345666", "text":text})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).Should(ContainSubstring("text is invalid"))
			})

			It("should get success", func() {
				url := fmt.Sprintf("/outbound/sms")
				text := "test2"
				request := preparePostRequest(url, bytes.NewBuffer(toJsonStr(map[string]string{"from":mockData.Phone2.Number, "to": mockData.Phone1.Number, "text":text})))
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).Should(ContainSubstring("outbound sms ok"))
			})

		})
	})
})
