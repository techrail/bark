package controllers

import (
	`encoding/json`
	`time`

	`github.com/valyala/fasthttp`

	`github.com/techrail/bark/constants`
)

type indexInfo struct {
	AppName        string `json:"appName"`
	AppVersion     string `json:"appVersion"`
	CurrentUtcTime string `json:"currentUtcTime"`
}

func IndexController(ctx *fasthttp.RequestCtx) {
	i := indexInfo{
		AppName:        constants.AppName,
		AppVersion:     constants.AppVersion,
		CurrentUtcTime: time.Now().UTC().Format(time.RFC3339),
	}

	iJson, err := json.Marshal(i)
	if err != nil {
		_, _ = ctx.WriteString("Welcome to Bark! If you are seeing this message, please contact the site admin.")
		return
	}
	ctx.Response.Header.Add(fasthttp.HeaderContentType, "application/json; charset=utf-8")
	ctx.Response.SetBodyString(string(iJson))
}
