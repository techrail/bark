package client

import (
	"context"
	"fmt"
	"github.com/techrail/bark/client/barkslogger"
	"github.com/techrail/bark/client/services/clientLogSender"
	"github.com/techrail/bark/client/services/ingestion"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/models"
)

// type webhook func(models.BarkLog) error

type Config struct {
	BaseUrl     string
	ErrorLevel  string
	ServiceName string
	SessionName string
	Slogger     *slog.Logger
	// AlertWebhook webhook
}

func (c *Config) parseMessage(msg string) models.BarkLog {
	l := models.BarkLog{
		ServiceName: c.ServiceName,
		SessionName: c.SessionName,
	}

	if len(msg) < 6 {
		l.Message = msg
		l.Code = constants.DefaultLogCode
		l.LogLevel = constants.DefaultLogLevel
		return l
	}

	// Look for `-` in the message
	pos := strings.Index(msg, "-")
	if pos < 1 {
		// There is no `-` in the message.
		l.Message = msg
		l.Code = constants.DefaultLogCode
		l.LogLevel = constants.DefaultLogLevel
		return l
	}

	if pos > len(msg)-3 {
		// There is no `-` in the message in any meaningful way
		l.Message = msg
		l.Code = constants.DefaultLogCode
		l.LogLevel = constants.DefaultLogLevel
		return l
	}

	// separate the message and meta info
	l.Message = strings.TrimSpace(msg[pos+1:])
	meta := strings.TrimSpace(msg[:pos])

	// Separate the code and level
	metas := strings.Split(meta, "#")
	if len(metas) > 2 {
		// Improperly formatted message
		l.Message = msg
		l.Code = constants.DefaultLogCode
		l.LogLevel = constants.DefaultLogLevel
		return l
	}

	if len(metas) == 1 {
		if len(metas[0]) > 16 {
			// Our code field is only 16 characters wide.
			l.Message = msg
			l.Code = constants.DefaultLogCode
			l.LogLevel = constants.DefaultLogLevel
			return l
		} else {
			l.Code = metas[0]
			l.LogLevel = constants.DefaultLogLevel
			return l
		}
	}

	if len(metas) == 2 {
		logLvl := strings.TrimSpace(metas[0])
		logCode := strings.TrimSpace(metas[1])

		if len(logLvl) != 1 || len(logCode) > constants.MaxLogCodelength || len(logCode) == 0 {
			// incorrectly formatted message
			l.Message = msg
			l.Code = constants.DefaultLogCode
			l.LogLevel = constants.DefaultLogLevel
			return l
		}

		l.LogLevel = getLogLevelFromCharacter(logLvl)
		l.Code = logCode

		//fmt.Println("-----------------", logLvl, "<>", logCode)
	}

	return l
}

func getLogLevelFromCharacter(s string) string {
	switch strings.ToUpper(s) {
	case "P":
		return constants.Panic
	case "A":
		return constants.Alert
	case "E":
		return constants.Error
	case "W":
		return constants.Warning
	case "N":
		return constants.Notice
	case "I":
		return constants.Info
	case "D":
		return constants.Debug
	default:
		return constants.Info
	}
}

func (c *Config) Panic(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Panic
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), barkslogger.LvlPanic, message)
	}
}
func (c *Config) Alert(message string) {
	// Todo: handle the alert webhook call here
	l := c.parseMessage(message)
	l.LogLevel = constants.Alert
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), barkslogger.LvlAlert, message)
	}
}
func (c *Config) Error(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Error
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Error(message)
	}
}
func (c *Config) Warn(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Warning
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Warn(message)
	}
}
func (c *Config) Notice(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Notice
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), barkslogger.LvlNotice, message)
	}
}
func (c *Config) Info(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Info
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Info(message)
	}
}
func (c *Config) Debug(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Debug
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Debug(message)
	}
}
func (c *Config) Println(message string) {
	l := c.parseMessage(message)
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Info(message)
	} else {
		// In addition to sending the log to server, we should also print it!
		fmt.Println(message)
	}
}

func (c *Config) Panicf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Panic
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), barkslogger.LvlPanic, message)
	}
}
func (c *Config) Alertf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Alert
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), barkslogger.LvlAlert, message)
	}
}
func (c *Config) Errorf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Error
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Error(message)
	}
}
func (c *Config) Warnf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Warning
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Warn(message)
	}
}
func (c *Config) Noticef(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Notice
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), barkslogger.LvlNotice, message)
	}
}
func (c *Config) Infof(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Info
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Info(message)
	}
}
func (c *Config) Debugf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Debug
	go ingestion.InsertSingleRequest(l)

	if c.Slogger != nil {
		c.Slogger.Debug(message)
	}
}

// func (c *Config) SetAlertWebhook(f webhook) {
// 	c.AlertWebhook = f
// }

func NewClient(url, errLevel, svcName, sessName string) *Config {
	if strings.TrimSpace(sessName) == "" {
		sessName = appRuntime.SessionName
		fmt.Printf("L#1L3WBF - Using %v as Session Name", sessName)
	}

	go clientLogSender.StartSendingLogs(url)

	return &Config{
		BaseUrl:     url,
		ErrorLevel:  errLevel,
		ServiceName: svcName,
		SessionName: sessName,
		Slogger:     barkslogger.New(os.Stdout),
	}
}

// WithCustomOut allows users to set output to custom writer instead of the default standard output
func (c *Config) WithCustomOut(out io.Writer) {
	c.Slogger = barkslogger.New(out)
}

// WithSlogHandler allows users to specify their own slog handler
func (c *Config) WithSlogHandler(handler slog.Handler) {
	c.Slogger = barkslogger.NewWithCustomHandler(handler)
}
