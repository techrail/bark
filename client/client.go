package client

import (
	"context"
	"fmt"
	"github.com/techrail/bark/typs/jsonObject"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/models"
)

type webhook func(models.BarkLog) error

type Config struct {
	BaseUrl      string
	ErrorLevel   string
	ServiceName  string
	SessionName  string
	BulkSend     bool
	Slogger      *slog.Logger
	AlertWebhook webhook
}

func (c *Config) parseMessage(msg string) models.BarkLog {
	l := models.BarkLog{
		ServiceName: c.ServiceName,
		SessionName: c.SessionName,
	}

	if len(msg) < 6 {
		l.Message = msg
		l.Code = constants.DefaultLogCode
		l.LogLevel = c.ErrorLevel
		return l
	}

	// Look for `-` in the message
	pos := strings.Index(msg, "-")
	if pos < 1 {
		// There is no `-` in the message.
		l.Message = msg
		l.Code = constants.DefaultLogCode
		l.LogLevel = c.ErrorLevel
		return l
	}

	if pos > len(msg)-3 {
		// There is no `-` in the message in any meaningful way
		l.Message = msg
		l.Code = constants.DefaultLogCode
		l.LogLevel = c.ErrorLevel
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
		l.LogLevel = c.ErrorLevel
		return l
	}

	if len(metas) == 1 {
		if len(metas[0]) > constants.MaxLogCodelength {
			// Our code field is only 16 characters wide.
			l.Message = msg
			l.Code = constants.DefaultLogCode
			l.LogLevel = c.ErrorLevel
			return l
		} else {
			l.Code = metas[0]
			l.LogLevel = c.ErrorLevel
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
			l.LogLevel = c.ErrorLevel
			return l
		}

		l.LogLevel = c.getLogLevelFromCharacter(logLvl)
		l.Code = logCode

		//fmt.Println("-----------------", logLvl, "<>", logCode)
	}

	return l
}

func (c *Config) getLogLevelFromCharacter(s string) string {
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
		return c.ErrorLevel
	}
}

func (c *Config) getCharacterFromLogLevel(logLevel string) string {
	switch logLevel {
	case constants.Panic:
		return "P"
	case constants.Alert:
		return "A"
	case constants.Error:
		return "E"
	case constants.Warning:
		return "W"
	case constants.Notice:
		return "N"
	case constants.Info:
		return "I"
	case constants.Debug:
		return "D"
	default:
		return "I"
	}
}

func (c *Config) dispatchLogMessage(l models.BarkLog) {
	if c.BulkSend {
		go InsertSingleRequest(l)
	} else {
		Wg.Add(1)
		go func() {
			defer Wg.Done()
			_, err := PostLog(c.BaseUrl+constants.SingleInsertUrl, l)
			if err.Severity == 1 {
				fmt.Println(err.Msg)
			}
		}()
	}
}

func (c *Config) Panic(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Panic
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), LvlPanic, message)
	}
}

func (c *Config) Alert(message string, blocking bool) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Alert
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.AlertWebhook != nil {
		if blocking {
			err := c.AlertWebhook(l)
			if err != nil {
				if c.Slogger != nil {
					c.Slogger.Log(context.Background(), LvlAlert, "unable to send alert")
				} else {
					fmt.Printf("E#1LR1V1 - Webhook failed to send. Error: %v | Original Log Message: %v\n", err, message)
				}
			}
		} else {
			go func() {
				err := c.AlertWebhook(l)
				if err != nil {
					if c.Slogger != nil {
						c.Slogger.Log(context.Background(), LvlAlert, "unable to send alert")
					} else {
						fmt.Printf("E#1LR1V1 - Webhook failed to send. Error: %v | Original Log Message: %v\n", err, message)
					}
				}
			}()
		}
	}

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), LvlAlert, message)
	}
}

func (c *Config) Error(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Error
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Error(message)
	}
}
func (c *Config) Warn(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Warning
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Warn(message)
	}
}
func (c *Config) Notice(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Notice
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), LvlNotice, message)
	}
}
func (c *Config) Info(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Info
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Info(message)
	}
}

func (c *Config) Debug(message string) {
	l := c.parseMessage(message)
	l.LogLevel = constants.Debug
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Debug(message)
	}
}

func (c *Config) Default(message string) {
	l := c.parseMessage(message)
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		switch l.LogLevel {
		case PANIC:
			c.Slogger.Log(context.Background(), LvlPanic, message)
		case ALERT:
			c.Slogger.Log(context.Background(), LvlAlert, message)
		case ERROR:
			c.Slogger.Error(message)
		case WARNING:
			c.Slogger.Warn(message)
		case NOTICE:
			c.Slogger.Log(context.Background(), LvlNotice, message)
		case DEBUG:
			c.Slogger.Debug(message)
		case INFO:
			fallthrough
		default:
			c.Slogger.Info(message)
		}
	}
}

func (c *Config) Raw(rawLog RawLog, returnError bool) error {
	// Try to parse the more data field
	moreData, err := jsonObject.ToJsonObject(rawLog.MoreData)
	if err != nil {
		// Cannot convert the contents of the MoreData field to JSON
		if returnError {
			return fmt.Errorf("E#1LSV6K - Could not parse moreData field as valid json")
		} else {
			// We will save the error
			moreData = jsonObject.EmptyNotNullJsonObject()
			_ = moreData.SetNewTopLevelElement(constants.MoreDataClientParseErrorMessage, err.Error())
			return nil
		}
	}

	l := models.BarkLog{
		LogTime:     rawLog.LogTime,
		LogLevel:    rawLog.LogLevel,
		ServiceName: rawLog.ServiceName,
		SessionName: rawLog.SessionName,
		Code:        rawLog.Code,
		Message:     rawLog.Message,
		MoreData:    moreData,
	}

	c.dispatchLogMessage(l)

	message := fmt.Sprintf("%v#%v - %v", c.getCharacterFromLogLevel(l.LogLevel), l.Code, l.Message)

	if c.Slogger != nil {
		switch l.LogLevel {
		case PANIC:
			c.Slogger.Log(context.Background(), LvlPanic, message)
		case ALERT:
			c.Slogger.Log(context.Background(), LvlAlert, message)
		case ERROR:
			c.Slogger.Error(message)
		case WARNING:
			c.Slogger.Warn(message)
		case NOTICE:
			c.Slogger.Log(context.Background(), LvlNotice, message)
		case DEBUG:
			c.Slogger.Debug(message)
		case INFO:
			fallthrough
		default:
			c.Slogger.Info(message)
		}
	}
	return nil
}

func (c *Config) Println(message string) {
	l := c.parseMessage(message)
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		switch l.LogLevel {
		case PANIC:
			c.Slogger.Log(context.Background(), LvlPanic, message)
		case ALERT:
			c.Slogger.Log(context.Background(), LvlAlert, message)
		case ERROR:
			c.Slogger.Error(message)
		case WARNING:
			c.Slogger.Warn(message)
		case NOTICE:
			c.Slogger.Log(context.Background(), LvlNotice, message)
		case DEBUG:
			c.Slogger.Debug(message)
		case INFO:
			fallthrough
		default:
			c.Slogger.Info(message)
		}
	} else {
		// In addition to sending the log to server, we should also print it!
		fmt.Println(message)
	}
}

func (c *Config) Printf(message string, format ...any) {
	msg := fmt.Sprintf(message, format...)
	l := c.parseMessage(msg)
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		switch l.LogLevel {
		case PANIC:
			c.Slogger.Log(context.Background(), LvlPanic, msg)
		case ALERT:
			c.Slogger.Log(context.Background(), LvlAlert, msg)
		case ERROR:
			c.Slogger.Error(msg)
		case WARNING:
			c.Slogger.Warn(msg)
		case NOTICE:
			c.Slogger.Log(context.Background(), LvlNotice, msg)
		case DEBUG:
			c.Slogger.Debug(msg)
		case INFO:
			fallthrough
		default:
			c.Slogger.Info(msg)
		}
	} else {
		// In addition to sending the log to server, we should also print it!
		fmt.Println(msg)
	}
}

func (c *Config) Panicf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Panic
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), LvlPanic, message)
	}
}

func (c *Config) Alertf(message string, blocking bool, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Alert
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.AlertWebhook != nil {
		if blocking {
			err := c.AlertWebhook(l)
			if err != nil {
				if c.Slogger != nil {
					c.Slogger.Log(context.Background(), LvlAlert, "unable to send alert")
				} else {
					fmt.Printf("E#1LR1V1 - Webhook failed to send. Error: %v | Original Log Message: %v\n", err, message)
				}
			}
		} else {
			go func() {
				err := c.AlertWebhook(l)
				if err != nil {
					if c.Slogger != nil {
						c.Slogger.Log(context.Background(), LvlAlert, "unable to send alert")
					} else {
						fmt.Printf("E#1LR1V1 - Webhook failed to send. Error: %v | Original Log Message: %v\n", err, message)
					}
				}
			}()
		}
	}

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), LvlAlert, message)
	}
}

func (c *Config) Errorf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Error
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Error(message)
	}
}

func (c *Config) Warnf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Warning
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Warn(message)
	}
}

func (c *Config) Noticef(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Notice
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Log(context.Background(), LvlNotice, message)
	}
}

func (c *Config) Infof(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Info
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Info(message)
	}
}

func (c *Config) Debugf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	l := c.parseMessage(message)
	l.LogLevel = constants.Debug
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Debug(message)
	}
}

func (c *Config) SetAlertWebhook(f webhook) {
	c.AlertWebhook = f
}

func NewClient(url, errLevel, svcName, sessName string, enableSlog bool, enableBulkSend bool) *Config {
	if !isValid(errLevel) {
		fmt.Printf("L#1LPYG2 - %v is not an acceptable log level. %v will be used as the default log level", errLevel, constants.DefaultLogLevel)
		errLevel = INFO
	}

	if strings.TrimSpace(svcName) == "" {
		sessName = constants.DefaultLogServiceName
		fmt.Printf("L#1L3WBF - Blank service name supplied. Using %v as Session Name", sessName)
	}

	if strings.TrimSpace(sessName) == "" {
		sessName = appRuntime.SessionName
		fmt.Printf("L#1L3WBF - Blank session name supplied. Using %v as Session Name", sessName)
	}

	//Wg.Add(1)

	if enableBulkSend {
		go keepSendingLogs(url)
	}

	var slogger *slog.Logger

	if enableSlog {
		slogger = New(os.Stdout)
	} else {
		slogger = nil
	}

	return &Config{
		BaseUrl:     url,
		ErrorLevel:  errLevel,
		ServiceName: svcName,
		SessionName: sessName,
		Slogger:     slogger,
		BulkSend:    enableBulkSend,
	}
}

// WithCustomOut allows users to set output to custom writer instead of the default standard output
func (c *Config) WithCustomOut(out io.Writer) {
	c.Slogger = New(out)
}

// WithSlogHandler allows users to specify their own slog handler
func (c *Config) WithSlogHandler(handler slog.Handler) {
	c.Slogger = NewWithCustomHandler(handler)
}

func (c *Config) WaitAndEnd() {
	Wg.Wait()
}
