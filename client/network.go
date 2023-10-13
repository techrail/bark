package client

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"

	"github.com/techrail/bark/models"
	"github.com/techrail/bark/typs/appError"
)

// Todo: Write Issue: Bark isn't throwing an error when insertion fails.

func post(url, payload string) (string, appError.AppErr) {
	var appErr appError.AppErr
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBodyString(payload)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)

	if err != nil {
		appErr.Msg = fmt.Sprintf("Error when making network call: %v", err)
		appErr.Code = "E#1LUJGZ"
		appErr.Severity = 1

		return "", appErr
	}

	bodyBytes := resp.Body()

	if resp.Header.StatusCode() != fasthttp.StatusOK {
		appErr.Msg = fmt.Sprintf("POST request failed. Code: %v | Message: %v", resp.Header.StatusCode(), string(resp.Body()))
		appErr.Code = "E#1L3T9W"
		appErr.Severity = 1
	}

	return string(bodyBytes), appErr
}

func PostLog(url string, log models.BarkLog) (string, appError.AppErr) {
	logRawJson, _ := json.Marshal(log)
	return post(url, string(logRawJson))
}

func PostLogArray(url string, log []models.BarkLog) (string, appError.AppErr) {
	logRawJson, _ := json.Marshal(log)
	return post(url, string(logRawJson))
}

func Get(url string) (string, appError.AppErr) {
	var err appError.AppErr
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)

	bodyBytes := resp.Body()

	if resp.Header.StatusCode() != fasthttp.StatusOK {
		err.Msg = "Get request failed"
		err.Code = "E#U4N3ER"
		err.Severity = 1
	}

	return string(bodyBytes), appError.AppErr{}
}
