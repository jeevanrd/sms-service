package smsservice

import (
	"context"
	"github.com/golang/glog"
	"github.com/jeevanrd/sms-service/database"
	"time"
	"strings"
	"errors"
	"github.com/jeevanrd/sms-service/statusErrors"
	"github.com/jeevanrd/sms-service/utils"
	"net/http"
)

type SmsService interface {
	InboundSms(ctx context.Context, req SmsRequest) (interface{},error)
	OutboundSms(ctx context.Context, req SmsRequest) (interface{},error)
}

type smsService struct {
	db database.Repository
	cache utils.Cache
}

type InboundSmsResponse struct {

}

type OutboundSmsResponse struct {

}

type  Response struct {
	Message		string `json:"message"`
	Error		string `json:"error"`
}

func (s *smsService) InboundSms(ctx context.Context, req SmsRequest) (interface{},error) {
	//get account
	account := ctx.Value("account").(database.Account)
	_,err := s.db.GetPhoneNumber(account.Id, req.To)

	if (err != nil) {
		return InboundSmsResponse{}, statusErrors.New(errors.New("to parameter not found"), http.StatusNotFound)
	}

	text := replaceNewLines(req.Text);
	trimmedText := strings.TrimSpace(text);
	cachekey := req.From + "_stop_" + req.To

	textInCache,err := s.cache.GetIntValueFromCache(cachekey)

	if (textInCache == 0 && strings.ToLower(trimmedText) == "stop") {
		s.cache.SetIntValueInCache(cachekey, 1, 4 * time.Hour)
	}
	return Response{"inbound sms ok", ""}, nil
}

func (s *smsService) OutboundSms(ctx context.Context, req SmsRequest) (interface{},error){
	//get account
	account := ctx.Value("account").(database.Account)
	_,err := s.db.GetPhoneNumber(account.Id, req.From)

	if(err != nil) {
		return InboundSmsResponse{}, statusErrors.New(errors.New("from parameter not found"), http.StatusNotFound)
	}

	toFromPair := req.To + "_stop_" + req.From
	fromToPair := req.From + "_stop_" + req.To

	val,err := s.cache.GetIntValueFromCache(toFromPair)

	if (val > 0) {
		return InboundSmsResponse{}, statusErrors.New(errors.New("sms from " + req.From + " to " + req.To  + " blocked by STOP request"), http.StatusBadRequest)
	}

	val,err = s.cache.GetIntValueFromCache(fromToPair)

	if(val > 0) {
		return InboundSmsResponse{}, statusErrors.New(errors.New("sms from " + req.From + " to " + req.To  + " blocked by STOP request"), http.StatusBadRequest)
	}

	key := "from_" + req.From;
	count,err := s.cache.GetIntValueFromCache(key)

	if(count > 0) {
		if(count >= 50) {
			return InboundSmsResponse{}, statusErrors.New(errors.New("limit reached for from " + req.From), http.StatusTooManyRequests)
		}
		s.cache.Increment(key, 1);
	} else {
		s.cache.SetIntValueInCache(key, 1, 24 * time.Hour)
	}

	return Response{"outbound sms ok", ""}, nil
}

func replaceNewLines(val string) string {
	strings.Replace(val, "\r ", "", -1)
	strings.Replace(val, "\n ", "", -1)
	return val
}

//NewService creates new InstanceService
func NewService(repo database.Repository, cache utils.Cache) SmsService {
	glog.Info("Creating a new instance service with:", "Repository", repo)
	return &smsService{db: repo, cache:cache}
}