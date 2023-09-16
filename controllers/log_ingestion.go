package controllers

import (
	`encoding/json`

	`github.com/valyala/fasthttp`

	`github.com/techrail/bark/appRuntime`
	`github.com/techrail/bark/models`
	`github.com/techrail/bark/services/ingestion`
)

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
		ctx.Error("E#1KDWRF - Invalid request body structure", fasthttp.StatusBadRequest)
		return
	}

	go ingestion.InsertMultiple(multipleLogEntries)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func ShutdownService(ctx *fasthttp.RequestCtx) {
	appRuntime.ShutdownRequested.Store(true)
}
