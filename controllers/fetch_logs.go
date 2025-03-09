package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/techrail/bark/models"
	"github.com/valyala/fasthttp"
)

func FetchLogs(ctx *fasthttp.RequestCtx) {
	logLevel := string(ctx.QueryArgs().Peek("logLevel"))
	serviceName := string(ctx.QueryArgs().Peek("serviceName"))
	sessionName := string(ctx.QueryArgs().Peek("sessionName"))
	startDate := string(ctx.QueryArgs().Peek("startDate"))
	endDate := string(ctx.QueryArgs().Peek("endDate"))
	barkLogDao := models.NewBarkLogDao()
	logs, err := barkLogDao.FetchLogs(logLevel, serviceName, sessionName, startDate, endDate)
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBodyString(fmt.Sprintf("Error fetching logs: %s", err.Error()))
		return
	}
	jsonResponse, err := json.Marshal(logs)
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBodyString("Error processing logs.")
		return
	}
	ctx.Response.SetBody(jsonResponse)
}
