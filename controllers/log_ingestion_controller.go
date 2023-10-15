package controllers

import (
	"encoding/json"

	"github.com/valyala/fasthttp"

	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/models"
	"github.com/techrail/bark/services/ingestion"
)

// SendSingleToChannel is tasked with handling all the single log insertion requests.
// It will simply send a response code 503 (service unavailable) if server shut down has already been requested.
// Response code 400 will be returned in case request body is empty.
// It will unmarshal the request body and compare the structure with structure of BarkLog struct.
// Finally, it will spawn a go routine to send the log to LogChannel and will respond with 200 to the client.
func SendSingleToChannel(ctx *fasthttp.RequestCtx) {
	if appRuntime.ShutdownRequested.Load() == true {
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return
	}

	body := ctx.Request.Body()
	if len(body) == 0 {
		ctx.Error("E#1KDWEO - Empty request", fasthttp.StatusBadRequest)
		return
	}

	var singleLogEntry models.BarkLog
	if err := json.Unmarshal(body, &singleLogEntry); err != nil {
		ctx.Error("E#1KDWET - Invalid request body structure: "+err.Error(), fasthttp.StatusBadRequest)
		return
	}

	go ingestion.InsertSingle(singleLogEntry)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// SendMultipleToChannel is tasked with handling all the multiple log insertion requests.
// It will simply send a response code 503 (service unavailable) if server shut down has already been requested.
// Response code 400 will be returned in case request body is empty.
// It will unmarshal the request body and compare the structure with structure of BarkLog struct slice.
// Finally, it will spawn a go routine to send the logs to LogChannel and will respond with 200 to the client.
func SendMultipleToChannel(ctx *fasthttp.RequestCtx) {
	if appRuntime.ShutdownRequested.Load() == true {
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return
	}

	body := ctx.Request.Body()
	if len(body) == 0 {
		ctx.Error("E#1KDWRA - Empty request", fasthttp.StatusBadRequest)
		return
	}

	var multipleLogEntries []models.BarkLog
	if err := json.Unmarshal(body, &multipleLogEntries); err != nil {
		ctx.Error("E#1KDWRF - Invalid request body structure: "+err.Error(), fasthttp.StatusBadRequest)
		return
	}

	go ingestion.InsertMultiple(multipleLogEntries)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// ShutdownService will set the value of global variable `ShutdownRequested` to true.
func ShutdownService(ctx *fasthttp.RequestCtx) {
	appRuntime.ShutdownRequested.Store(true)
}
