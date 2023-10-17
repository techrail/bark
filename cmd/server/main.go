package main

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/techrail/bark/models"
	"github.com/techrail/bark/utils"
	"github.com/valyala/fasthttp"
	"os"

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

	// The index endpoint displays information about the bark server.
	r.GET("/", controllers.IndexController)

	// This is a demo endpoint to ensure bark server is running which will print Hello, `name`!.
	r.GET("/hello/{name}", Hello)

	// Bark client contains the logic which decides which out of the two (single/multiple) insertion endpoints is called.
	// This endpoint sends single log entry at a time to the DB.
	r.POST("/insertSingle", controllers.SendSingleToChannel)
	// This endpoint handles the batch insertion of logs to the DB.
	r.POST("/insertMultiple", controllers.SendMultipleToChannel)

	// This endpoint is responsible to initiate server shut down, following which Bark server will not process any new incoming requests.
	// It will, however, shut down after it has completely saved all the logs received up till that point to the database.
	r.POST("/shutdownServiceAsap", controllers.ShutdownService)

	fmt.Printf("I#1M2UDR - Database connection string from Environment: %s\n", os.Getenv("BARK_DATABASE_URL"))

	dbUrl := os.Getenv("BARK_DATABASE_URL")
	err := utils.ParsePostgresUrl(dbUrl)
	if err != nil {
		panic(err.Error())
	}
	err = resources.InitDb(dbUrl)
	if err != nil {
		panic("E#1KDZRP - " + err.Error())
	}
	bld := models.NewBarkLogDao()

	// Sends a single log entry to the postgres DB stating Bark server has started successfully.
	// Returns an error and halts the server boot up in case the connection acquired to the postgres DB is not proper.
	err = bld.InsertServerStartedLog()
	if err != nil {
		panic("P#1LQ2YQ - Bark server start failed: " + err.Error())
	}

	// Go routine which writes logs received in the LogChannel to DB.
	go dbLogWriter.KeepSavingLogs()
	err = fasthttp.ListenAndServe(address, r.Handler)
	if err != nil {
		fmt.Println("E#1M30BA - Listen and serve failed: ", err.Error())
	}
}
