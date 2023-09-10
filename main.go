package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/techrail/bark/db"
	"github.com/techrail/bark/logger"
	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}
func Init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("ABCD" + os.Getenv("DB_PASSWORD"))
	fmt.Println(gofakeit.UUID())
}

func main() {
	Init()
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)

	// Create DB Connection
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully connected to database")

	// Test transactions
	sampleLog := logger.Log{
		// Id:          1234,
		LogTime:     time.Now(),
		LogLevel:    0,
		ServiceName: "test",
		Code:        "1234",
		Message:     "Test",
		MoreData:    "{}",
	}
	err = db.InsertLog(sampleLog)
	if err != nil {
		log.Fatal(err)
	}
	// log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
