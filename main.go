package main

import (
	"fmt"
	"log"
	"os"
	//"flag"
	//"strconv"

	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"

	"github.com/techrail/bark/controllers"
	"github.com/techrail/bark/resources"
	"github.com/techrail/bark/services/dbLogWriter"
)

//var port = flag.Uint("p",8080,"Enter the port number")
func Index(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("Welcome to Bark!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	_, _ = fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

// Init performs prerequisite tasks - like loading env variables
func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//flag.Parse()
	//p := ":"+strconv.FormatUint(uint64(*port),10)
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)
	r.POST("/insertSingle", controllers.SendSingleToChannel)
	r.POST("/insertMultiple", controllers.SendMultipleToChannel)
	r.POST("/shutdownServiceAsap", controllers.ShutdownService)

	err := resources.InitDB()
	if err != nil {
		log.Fatal("E#1KDZRP - " + err.Error())
	}
	go dbLogWriter.StartWritingLogs()
	log.Fatal(fasthttp.ListenAndServe(os.Getenv("APP_PORT"), r.Handler))
}
