package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/fasthttp/router"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/techrail/bark/utils"
	"github.com/valyala/fasthttp"
)

type LogLevel string

func (l LogLevel) ToInt() int16 {
	switch l {
	case "INFO":
		return 10
	case "DEBUG":
		return 11
	case "WARNING":
		return 12
	default:
		return 10
	}
}

type Log struct {
	LogLevel LogLevel               `json:"log_level"`
	SVCName  string                 `json:"service_name"`
	Code     string                 `json:"code"`
	Msg      string                 `json:"msg"`
	MoreData map[string]interface{} `json:"more_data"`
}

type App struct {
	DB *sqlx.DB
}
type AppConfig struct {
	DBName     string
	DBUsername string
	DBPassword string
	SSLMode    string
	APPPort    string
}

func NewApp(opts utils.AppConfig) (*App, error) {
	dSN := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", opts.DBUsername, opts.DBPassword, opts.DBName, opts.SSLMode)
	db, err := sqlx.Connect("postgres", dSN)

	if err != nil {
		return nil, err
	}

	return &App{
		DB: db,
	}, nil
}

func buildInserQuery(l Log) (string, []any) {
	var queryBuilder, valuesBuilder strings.Builder
	var elems []any

	queryBuilder.WriteString("INSERT INTO app_log (")
	valuesBuilder.WriteString("VALUES (")

	v := reflect.ValueOf(l)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i).Interface()

		if reflect.ValueOf(value).IsZero() {
			continue
		}
		queryBuilder.WriteString(field.Tag.Get("json") + ",")
		valuesBuilder.WriteString(fmt.Sprintf("$%d,", len(elems)+1))

		switch val := value.(type) {
		case LogLevel:
			elems = append(elems, val.ToInt())
		case string:
			elems = append(elems, val)
		case map[string]interface{}:
			mD, err := json.Marshal(val)
			if err != nil {
				continue
			}

			elems = append(elems, mD)
		}
	}
	query := queryBuilder.String()[:queryBuilder.Len()-1] + ") " + valuesBuilder.String()[:valuesBuilder.Len()-1] + ")"

	return query, elems
}

func (a App) Insert(ctx *fasthttp.RequestCtx) {
	b := ctx.PostBody()
	var l Log
	err := json.Unmarshal(b, &l)

	if err != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBodyString(err.Error())
		return
	}

	tx, err := a.DB.Begin()
	if err != nil {
		ctx.Response.SetStatusCode(500)
		log.Println(err)
		ctx.Response.SetBodyString("failed to insert")
		return
	}

	query, elems := buildInserQuery(l)
	if len(elems) == 0 {
		ctx.Response.SetStatusCode(400)
		ctx.Response.SetBodyString("empty request")
		return
	}
	_, err = tx.Exec(query, elems...)

	if err != nil {
		log.Println(err)
	}

	err = tx.Commit()
	if err != nil {
		ctx.Response.SetStatusCode(500)
		log.Println(err)
		ctx.Response.SetBodyString("failed to insert")
		return
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

}
