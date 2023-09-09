package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

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

type LogD struct {
	LogLevel int16
	SVCName  string
	Code     string
	Msg      string
	MoreData map[string]interface{}
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
	query := "INSERT INTO app_log ("
	values := "VALUES ("

	t := reflect.TypeOf(l)
	v := reflect.ValueOf(l)

	j := 1

	elems := []any{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		if reflect.ValueOf(value).IsZero() {
			continue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			var mD []byte
			// marshal moredata if it exists
			if len(l.MoreData) != 0 {
				mD, _ = json.Marshal(l.MoreData)
				elems = append(elems, mD)
			}

		case reflect.String:
			keyName := field.Tag.Get("json")
			if keyName == "log_level" {
				val := value.(LogLevel)
				elems = append(elems, val.ToInt())
			} else {
				elems = append(elems, value)
			}
		}

		// Append the field name to the query
		query += field.Tag.Get("json")

		// Append the placeholder to the values
		values += fmt.Sprintf("$%d", j)
		j++

		// Add commas between fields and placeholders
		if i < t.NumField()-1 {
			query += ", "
			values += ", "
		}
	}

	query += ") " + values + ")"
	log.Println(query)

	_, err = tx.Exec(query, elems...)

	if err != nil {
		fmt.Println(err)
	}

	err = tx.Commit()
	if err != nil {
		ctx.Response.SetStatusCode(500)
		log.Println(err)
		ctx.Response.SetBodyString("failed to insert")
		return
	}

	ctx.Response.SetBody(b)
	ctx.Response.SetStatusCode(200)
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
}
