package controllers

import (
	`encoding/json`

	`github.com/valyala/fasthttp`

	`github.com/techrail/bark/models`
	`github.com/techrail/bark/services/ingestion`
)

func SendSingleToChannel(ctx *fasthttp.RequestCtx) {
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
	// ctx.Response.SetBodyString("Sent for insertion")
}

func SendMultipleToChannel(ctx *fasthttp.RequestCtx) {
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
	// ctx.Response.SetBodyString("Sent for insertion")
}
