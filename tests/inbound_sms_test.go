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
)

func init() {
	go startServer()
}

func TestInboundSmsMethodContentType(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", nil)
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
	assert.Equal(t, result.Error, "invalid content type")
}

func TestInboundSmsNonPostMethods(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("GET", "http://localhost:7070/inbound/sms", nil)
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

func TestInboundSmsAuthFailure(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", nil)
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

func TestInboundSmsUnknownFailureScenarios(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", nil)
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
	assert.Equal(t, result.Error, "invalid payload")
}

func TestInboundSmsBadRequestScenarioswithFormParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{}`)))
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

	req, _ = http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "12345"}`)))
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


func TestInboundSmsBadRequestScenarioswithToParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232"}`)))
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

	req, _ = http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to":"12323"}`)))
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

func TestInboundSmsBadRequestScenarioswithTextParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to": "12323232"}`)))
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

	req, _ = http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to":"12323123", "text":""}`)))
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

func TestInboundSmsWithNotFoundFromParameter(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to":"12323123", "text":"testing"}`)))
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
	assert.Equal(t, result.Error, "to parameter not found")
}


func TestInboundSmsSucessScenario(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("POST", "http://localhost:7070/inbound/sms", bytes.NewBuffer([]byte(`{"from": "123451232232", "to":"4924195509198", "text":"testing"}`)))
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
}

