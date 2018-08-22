package smsservice_test

import (
	dbMocks "github.com/jeevanrd/sms-service/database/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"context"
	"github.com/jeevanrd/sms-service/mockData"
	"github.com/go-kit/kit/log"
	"os"
	"github.com/jeevanrd/sms-service/smsservice"
	"github.com/jeevanrd/sms-service/utils"
	"errors"
	"github.com/jeevanrd/sms-service/statusErrors"
	"net/http"
	"github.com/jeevanrd/sms-service/database"
)



var _ = Describe("Sms Service", func() {

	ctx := context.WithValue(context.Background(), "account", mockData.Account1)
	mockDatabaseRepo := &dbMocks.Repository{}

	mockDatabaseRepo.On("GetAccount", mockData.Account1.Username, mockData.Account1.AuthId).Return(mockData.Account1, nil)
	mockDatabaseRepo.On("GetPhoneNumber", mockData.Phone1.AccountId, mockData.Phone2.Number).Return(mockData.Phone2, nil)
	mockDatabaseRepo.On("GetPhoneNumber", mockData.Phone1.AccountId, mockData.Phone1.Number).Return(mockData.Phone1, nil)
	mockDatabaseRepo.On("GetPhoneNumber", mockData.Phone1.AccountId, "112345678").Return(database.PhoneNumber{}, errors.New("not found"))
	mockDatabaseRepo.On("GetPhoneNumber", mockData.Phone1.AccountId, "12345678").Return(database.PhoneNumber{}, errors.New("not found"))

	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var is smsservice.SmsService
	redisClient := utils.GetCacheClient()
	is = smsservice.NewService(mockDatabaseRepo, utils.LocalCache{Client:redisClient})

	Describe("GET /inbound/sms", func() {
		Context("test inbound", func() {
			It("error should be nil", func() {
				var sampleSmsR1 = smsservice.SmsRequest{From:mockData.Phone1.Number,To:mockData.Phone2.Number,Text:"text1"}
				output,err := is.InboundSms(ctx,sampleSmsR1)
				resp := output.(smsservice.Response)
				Expect(err).To(BeNil())
				Expect(resp.Message).To(Equal("inbound sms ok"))
			})

			It("should get not found", func() {
				var sampleSmsR1 = smsservice.SmsRequest{From:mockData.Phone1.Number,To:"12345678",Text:"test1"}
				_,err := is.InboundSms(ctx,sampleSmsR1)
				Expect(err).To(Equal(statusErrors.New(errors.New("to parameter not found"), http.StatusNotFound)))
			})
		})
	})


	Describe("GET /outbound/sms", func() {
		Context("test outbound sms", func() {
			It("error should be nil", func() {
				var sampleSmsR1 = smsservice.SmsRequest{From:mockData.Phone1.Number,To:mockData.Phone3.Number,Text:"test1"}
				output,err := is.OutboundSms(ctx,sampleSmsR1)
				resp := output.(smsservice.Response)
				Expect(err).To(BeNil())
				Expect(resp.Message).To(Equal("outbound sms ok"))
			})

			It("should get not found", func() {
				var sampleSmsR1 = smsservice.SmsRequest{From:"112345678",To:mockData.Phone1.Number,Text:"test1"}
				_,err := is.OutboundSms(ctx,sampleSmsR1)
				Expect(err).To(Equal(statusErrors.New(errors.New("from parameter not found"), http.StatusNotFound)))
			})


			It("should get blocked error", func() {
				var sampleSmsR1 = smsservice.SmsRequest{From:mockData.Phone1.Number,To:mockData.Phone2.Number,Text:"stop"}
				_,err := is.InboundSms(ctx,sampleSmsR1)
				_,err = is.OutboundSms(ctx,sampleSmsR1)
				Expect(err).To(Equal(statusErrors.New(errors.New("sms from 4924195509198 to 4924195509196 blocked by STOP request"), http.StatusBadRequest)))
			})
		})
	})
})
