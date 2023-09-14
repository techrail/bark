package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"

	`github.com/techrail/bark/controllers`
	`github.com/techrail/bark/resources`
)

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
	Init()
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)
	r.POST("/insertSingle", controllers.SendSingleToChannel)
	r.POST("/insertMultiple", controllers.SendMultipleToChannel)

	err := resources.InitDatabase()
	if err != nil {
		log.Fatal("E#1KDZRP - " + err.Error())
	}
	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))

	// =========== TEST CASE (To be refactored) ===========
	// NOTE: We will write the tests later, separately

	// more_data is a JSONB field in the db, in the BarkLog struct its stored as a json.RawMessage ([]byte) field.
	// So we need to Marshal it to json before inserting
	// moreData, _ := json.Marshal(map[string]interface{}{
	// 	"a": "apple",
	// 	"b": "banana",
	// },
	// )
	// sampleLog := []models.BarkLog{
	// 	// Id:          1234,
	// 	{LogTime: time.Now(),
	// 		LogLevel:    0,
	// 		ServiceName: "test",
	// 		Code:        "1234",
	// 		Message:     "Test",
	// 		MoreData:    moreData},
	// 	{LogTime: time.Now(),
	// 		LogLevel:    0,
	// 		ServiceName: "test",
	// 		Code:        "1234",
	// 		Message:     "Test",
	// 		MoreData:    moreData},
	// 	{LogTime: time.Now(),
	// 		LogLevel:    0,
	// 		ServiceName: "test",
	// 		Code:        "1234",
	// 		Message:     "Test",
	// 		MoreData:    moreData},
	// 	{LogTime: time.Now(),
	// 		LogLevel:    0,
	// 		ServiceName: "test",
	// 		Code:        "1234",
	// 		Message:     "Test",
	// 		MoreData:    moreData},
	// }
	// err = barkDb.InsertBatch(sampleLog)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // 2.Fetch n number of logs
	// logs, err := barkDb.FetchLimitedLogs(4)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Print(logs)

	// =============================================================

	//	ctx.Response.SetStatusCode(200)

	// go batchCommit(logChannel)
}
