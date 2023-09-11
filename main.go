package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
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

	ctx.Response.SetStatusCode(200)
}

func readFromChannel(logChannel <-chan Log) []Log {
	var logsFromChannel []Log
	for val := range logChannel {
		logsFromChannel = append(logsFromChannel, val)
	}

	return logsFromChannel
}

func writeToChannel(logChannel chan Log, logRecord Log) {
	logChannel <- logRecord
}

// go routine to check channel length and commit to DB
func batchCommit(logChannel chan Log) string {
	logChannelLength := 0
	for {
		logChannelLength = len(logChannel)
		if logChannelLength > 100 {
			//commit in batches of 100
			logsToCommit := readFromChannel(logChannel)
			for i := 0; i < len(logsToCommit); i++ {
				// call bulk insert function
			}

		} else if logChannelLength > 0 && logChannelLength < 100 {
			// commit one at a time
			logsToCommit := readFromChannel(logChannel)
			for i := 0; i < len(logsToCommit); i++ {
				// call insert
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

func main() {

	config := utils.LoadConfig()

	app, err := NewApp(config)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("connected to database")

	r := router.New()
	r.POST("/insert", app.Insert)

	log.Fatal(fasthttp.ListenAndServe(config.APPPort, r.Handler))

	logChannel := make(chan Log)

	go batchCommit(logChannel)

	for _, log := range logs {
		fmt.Printf("%v\n", log)
	}

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
