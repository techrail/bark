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

var logChannel = make(chan models.BarkLog)

func main() {
	Init()
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)
	r.POST("/insert", SendToChannel)

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
	sampleLog := []models.BarkLog{
		// Id:          1234,
		{LogTime: time.Now(),
			LogLevel:    0,
			ServiceName: "test",
			Code:        "1234",
			Message:     "Test",
			MoreData:    moreData},
		{LogTime: time.Now(),
			LogLevel:    0,
			ServiceName: "test",
			Code:        "1234",
			Message:     "Test",
			MoreData:    moreData},
		{LogTime: time.Now(),
			LogLevel:    0,
			ServiceName: "test",
			Code:        "1234",
			Message:     "Test",
			MoreData:    moreData},
		{LogTime: time.Now(),
			LogLevel:    0,
			ServiceName: "test",
			Code:        "1234",
			Message:     "Test",
			MoreData:    moreData},
	}
	err = db.InsertBatch(sampleLog)
	if err != nil {
		log.Fatal(err)
	}

	// 2.Fetch n number of logs
	logs, err := db.FetchLimitedLogs(4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(logs)

	//	ctx.Response.SetStatusCode(200)

	go batchCommit(logChannel)
}

func SendToChannel(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.Body()
	if len(body) == 0 {
		ctx.Error("Empty request", fasthttp.StatusBadRequest)
		return
	}
	var requestData models.BarkLog
	if err := json.Unmarshal(body, &requestData); err != nil {
		ctx.Error("Invalid request body structure", fasthttp.StatusBadRequest)
		return
	}

	select {
	case logChannel <- requestData:
	default:
		ctx.Error("Channel is busy", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

}

// go routine to check channel length and commit to DB
func batchCommit(logChannel chan models.BarkLog) string {
	logChannelLength := 0
	db := new(db.BarkPostgresDb)
	for {
		logChannelLength = len(logChannel)
		if logChannelLength > 100 {
			var logBatch = []models.BarkLog{}
			for i := 0; i < 100; i++ {
				elem, ok := <-logChannel
				if !ok {
					fmt.Println("Error occured while getting batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}
			err := db.InsertBatch(logBatch)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Batch inserted at ", time.Now().Format("2006-01-02 15:04:05"))

		} else if logChannelLength > 0 && logChannelLength < 100 {
			// commit one at a time
			singleLog := <-logChannel
			err := db.InsertLog(singleLog)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Log inserted at ", time.Now().Format("2006-01-02 15:04:05"))

		} else {
			time.Sleep(1 * time.Second)
		}
	}

}
