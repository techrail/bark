package client

import (
	"encoding/json"
	"github.com/techrail/bark/models"
	"github.com/techrail/bark/typs/appError"
	"github.com/valyala/fasthttp"
)

// Todo: Write Issue: Bark isn't throwing an error when insertion fails.

// Todo: Write Issue: Bark to send proper JSON response after insertion.

func post(url, payload string) (string, appError.AppErr) {
	var err appError.AppErr
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBodyString(payload)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	_ = client.Do(req, resp)

	bodyBytes := resp.Body()

	if resp.Header.StatusCode() != fasthttp.StatusOK {
		err.Msg = "POST request failed"
		err.Code = "E#3HMV3G"
		err.Severity = 1
	}

	return string(bodyBytes), err
}

func PostLog(url string, log models.BarkLog) (string, appError.AppErr) {
	logRawJson, _ := json.Marshal(log)
	return post(url, string(logRawJson))
}

func PostLogs(url string, log []models.BarkLog) (string, appError.AppErr) {
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
