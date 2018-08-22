package tests

import (
	"testing"
	"net/http"
	"time"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"github.com/jeevanrd/sms-service/auth"
	"github.com/gorilla/mux"
	"github.com/jeevanrd/sms-service/routers"
)

func init() {
	go startServer()
}

func startServer() {
	router := mux.NewRouter().StrictSlash(true)
	routers.CreateAppRouter(router)
	http.ListenAndServe(":7070", router)
}

func TestOutboundSmsMethodContentType(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", nil)
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 415, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "Unsupported Media Type")
}

func TestOutboundSmsNonPostMethods(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("GET", "http://localhost:7070/outbound/sms", nil)
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 405, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "Method Not Allowed")
}

func TestOutboundSmsAuthFailure(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", nil)
	req.Header.Add("content-type", "application/json")
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 403, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "Please pass valid credentials")
}

func TestOutboundSmsBadRequestScenarioswithFormParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "from is missing")

	req, _ = http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "12345"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp,err = client.Do(req)
	if err != nil {
		panic(err)
	}

	rawdata, err = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "from is invalid")
}


func TestOutboundSmsBadRequestScenarioswithToParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "to is missing")

	req, _ = http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to":"12323"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp,err = client.Do(req)
	if err != nil {
		panic(err)
	}

	rawdata, err = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "to is invalid")
}

func TestOutboundSmsBadRequestScenarioswithTextParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to": "12323232"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "text is missing")

	req, _ = http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to":"12323123", "text":""}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp,err = client.Do(req)
	if err != nil {
		panic(err)
	}

	rawdata, err = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "text is invalid")
}

func TestOutboundSmsWithNotFoundFromParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to":"12323123", "text":"testing"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 404, resp.StatusCode)
	assert.Equal(t, result.Message, "")
	assert.Equal(t, result.Error, "from parameter not found")
}


func TestOutboundSmsSucessScenario(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "4924195509197", "to":"3253280312", "text":"testing"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, result.Error, "")
	assert.Equal(t, result.Message, "outbound sms ok")
}


func TestOutboundSmsSTOPScenario(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "3253280312", "to":"4924195509198", "text":"stop"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, result.Error, "")
	assert.Equal(t, result.Message, "inbound sms ok")


	req, _ = http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "4924195509198", "to":"3253280312", "text":"testing"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	rawdata, err = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, result.Error, "sms from 4924195509198 to 3253280312 blocked by STOP request")
	assert.Equal(t, result.Message, "")
}


func TestOutboundSmsRateLimitExceededScenario(t *testing.T) {

	for i := 0; i < 50; i++ {
		client := &http.Client{
			Timeout: 1 * time.Second,
		}
		req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "4924195509195", "to":"3253280312", "text":"testing"}`)))
		req.Header.Add("content-type", "application/json")
		req.SetBasicAuth("plivo1", "20S0KPNOIM")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, 200, resp.StatusCode)
	}

	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://localhost:7070/outbound/sms", bytes.NewBuffer([]byte(`{"from": "4924195509195", "to":"3253280312", "text":"testing"}`)))
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("plivo1", "20S0KPNOIM")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var result auth.Response
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rawdata, &result); err != nil {
		panic(err)
	}

	assert.Equal(t, 429, resp.StatusCode)
	assert.Equal(t, result.Error, "limit reached for from 4924195509195")
	assert.Equal(t, result.Message, "")

}
