package main

import (
	"fmt"
	"github.com/techrail/bark/models"
	"log"
	"os"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/techrail/bark/controllers"
	"github.com/techrail/bark/resources"
	"github.com/techrail/bark/services/dbLogWriter"
)

func Hello(ctx *fasthttp.RequestCtx) {
	_, _ = fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

// Init performs prerequisite tasks - like loading env variables
func Init() string {
	port := os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	address := ":" + port
	return address
}

func main() {
	address := Init()
	r := router.New()
	r.GET("/", controllers.IndexController)
	r.GET("/hello/{name}", Hello)
	r.POST("/insertSingle", controllers.SendSingleToChannel)
	r.POST("/insertMultiple", controllers.SendMultipleToChannel)
	r.POST("/shutdownServiceAsap", controllers.ShutdownService)
	err := resources.InitDb()
	if err != nil {
		log.Fatal("E#1KDZRP - " + err.Error())
	}
	bld := models.NewBarkLogDao()
	err = bld.InsertServerStartedLog()
	if err != nil {
		log.Fatal("P#1LQ2YQ - Bark server start failed: " + err.Error())
	}

	go dbLogWriter.StartWritingLogs()
	log.Fatal(fasthttp.ListenAndServe(address, r.Handler))

}
