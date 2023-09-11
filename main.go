package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/techrail/bark/db"
	"github.com/techrail/bark/models"
	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

// Perform prerequisite tasks - like loading env variables
func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	Init()
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)

	// Connect to Postgres DB instance
	db, err := db.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	// Ping DB
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully connected to database")

	// Test transactions
	// 1. Insert Log

	// more_data is a JSONB field in the db, in the BarkLog struct its stored as a json.RawMessage ([]byte) field.
	// So we need to Marshal it to json before inserting
	moreData, _ := json.Marshal(map[string]interface{}{
		"a": "apple",
		"b": "banana",
	},
	)
	sampleLog := models.BarkLog{
		// Id:          1234,
		LogTime:     time.Now(),
		LogLevel:    0,
		ServiceName: "test",
		Code:        "1234",
		Message:     "Test",
		MoreData:    moreData,
	}
	err = db.InsertLog(sampleLog)
	if err != nil {
		log.Fatal(err)
	}

	// 2.Fetch n number of logs
	logs, err := db.FetchLimitedLogs(4)
	if err != nil {
		log.Fatal(err)
	}
	for _, log := range logs {
		fmt.Printf("%v\n", log)
	}

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
