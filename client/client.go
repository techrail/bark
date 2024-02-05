package client

import (
	"context"
	"fmt"
	"github.com/techrail/bark/internal/jsonObject"
	"github.com/techrail/bark/resources"
	"github.com/techrail/bark/services/dbLogWriter"
	"github.com/techrail/bark/services/ingestion"
	"github.com/techrail/bark/utils"
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
	serverMode             int
	disableDebugLvlLogging bool
	BaseUrl                string
	ErrorLevel             string
	ServiceName            string
	ServiceInstanceName    string
	BulkSend               bool
	Slogger                *slog.Logger
	AlertWebhook           webhook
}

// parseMessage extracts LMID (Log Message Identifier) if a valid LMID exists in message string otherwise.
func (c *Config) parseMessage(msg string) models.BarkLog {
	l := models.BarkLog{
		ServiceName:         c.ServiceName,
		ServiceInstanceName: c.ServiceInstanceName,
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

// getLogLevelFromCharacter returns a string of log level for the first character of log level passed.
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

// getCharacterFromLogLevel returns a string of letter of log level for the log level string passed.
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

// dispatchLogMessage is responsible for dispatching a log to server.
// This method takes in a log and sends it to PendingLogsChan if bulk send is enabled,
// otherwise if bulk send is disabled it creates a network request to send the log,
// in a new goroutine.
func (c *Config) dispatchLogMessage(l models.BarkLog) {
	switch c.serverMode {
	case constants.ClientServerUsageModeDisabled:
		// Server is disabled. Nothing
		return
	case constants.ClientServerUsageModeRemoteServer:
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
	case constants.ClientServerUsageModeEmbedded:
		go ingestion.InsertSingle(l)
	default:
		panic(fmt.Sprintf("P#1M2RSP - Unexpected server usage mode: %v", c.serverMode))
	}
}

// Panic sends a LvlPanic log to server and prints the log if slog is enabled.
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

// Alert sends a LvlAlert log to server and prints the log if slog is enabled.
// This also initiates a AlertWebhook call.
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

// Error sends a LvlError level log to server and prints the log if slog is enabled.
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

// Warn sends a LvlWarning level log to server and prints the log if slog is enabled.
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

// Notice sends a LvlNotice level log to server and prints the log if slog is enabled.
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

// Info sends a LvlInfo level log to server and prints the log if slog is enabled.
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

// Debug sends a LvlDebug level log to server and prints the log if slog is enabled.
func (c *Config) Debug(message string) {
	if c.disableDebugLvlLogging {
		return
	}

	l := c.parseMessage(message)
	l.LogLevel = constants.Debug
	l.LogTime = time.Now().UTC()
	l.MoreData = jsonObject.EmptyNotNullJsonObject()
	c.dispatchLogMessage(l)

	if c.Slogger != nil {
		c.Slogger.Debug(message)
	}
}

// Default sends a Config.Default level log to server (specified while creating client)
// and prints the log if slog is enabled.
func (c *Config) Default(message string) {
	l := c.parseMessage(message)

	if c.disableDebugLvlLogging && l.LogLevel == DEBUG {
		return
	}

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

// Raw allows user to send a RawLog to server.
func (c *Config) Raw(rawLog RawLog, returnError bool) error {
	if c.disableDebugLvlLogging && rawLog.LogLevel == DEBUG {
		if returnError {
			return fmt.Errorf("E#1M7ITZ - Debug logs have been disabled. Log won't be processed")
		} else {
			// return silently
			return nil
		}
	}

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
		LogTime:             rawLog.LogTime,
		LogLevel:            rawLog.LogLevel,
		ServiceName:         rawLog.ServiceName,
		ServiceInstanceName: rawLog.ServiceInstanceName,
		Code:                rawLog.Code,
		Message:             rawLog.Message,
		MoreData:            moreData,
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

// Println sends logs to server based on LMID passed in message string.
// If the LMID is invalid an LvlInfo log level considered by default.
// This method prints the logs regardless if the slog is enabled or not.
func (c *Config) Println(message string) {
	l := c.parseMessage(message)

	if c.disableDebugLvlLogging && l.LogLevel == DEBUG {
		return
	}

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

// Printf performs the same operation as Config.Println but it accepts a format specifier.
func (c *Config) Printf(message string, format ...any) {
	msg := fmt.Sprintf(message, format...)
	c.Println(msg)
}

// Panicf performs the same operation as Config.Panic but it accepts a format specifier.
func (c *Config) Panicf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	c.Panic(message)
}

// Alertf performs the same operation as Config.Alert but it accepts a format specifier.
func (c *Config) Alertf(message string, blocking bool, format ...any) {
	message = fmt.Sprintf(message, format...)
	c.Alert(message, blocking)
}

// Errorf performs the same operation as Config.Error but it accepts a format specifier.
func (c *Config) Errorf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	c.Error(message)
}

// Warnf performs the same operation as Config.Warn but it accepts a format specifier.
func (c *Config) Warnf(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	c.Warn(message)
}

// Noticef performs the same operation as Config.Notice but it accepts a format specifier.
func (c *Config) Noticef(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	c.Notice(message)
}

// Infof performs the same operation as Config.Info but it accepts a format specifier.
func (c *Config) Infof(message string, format ...any) {
	message = fmt.Sprintf(message, format...)
	c.Info(message)
}

// Debugf performs the same operation as Config.Debug but it accepts a format specifier.
func (c *Config) Debugf(message string, format ...any) {
	if c.disableDebugLvlLogging {
		return
	}

	message = fmt.Sprintf(message, format...)
	c.Debug(message)
}

func (c *Config) DisableDebugLogs() {
	c.disableDebugLvlLogging = true
}

func (c *Config) EnableDebugLogs() {
	c.disableDebugLvlLogging = false
}

// SetAlertWebhook sets the function f to be used as a webhook to be used by Config.Alert method
func (c *Config) SetAlertWebhook(f webhook) {
	c.AlertWebhook = f
}

// NewSloggerClient creates and returns a new Config object which can be used for fully client-side logging only
// It accepts one single parameter - the default log level. The config returned cannot be used to send logs to any
// remote server. This can be used for projects that do not aim to send logs to a remote service yet!
func NewSloggerClient(defaultLogLevel string) *Config {
	if !isValid(defaultLogLevel) {
		fmt.Printf("L#1LZAY0 - %v is not an acceptable log level. %v will be used as the default log level", defaultLogLevel, constants.DefaultLogLevel)
		defaultLogLevel = constants.DefaultLogLevel
	}

	slogger := newSlogger(os.Stdout)

	return &Config{
		serverMode:          constants.ClientServerUsageModeDisabled,
		BaseUrl:             constants.DisabledServerUrl,
		ErrorLevel:          defaultLogLevel,
		ServiceName:         "",
		ServiceInstanceName: "",
		Slogger:             slogger,
		BulkSend:            false,
	}
}

// NewClient creates and returns a new Config object with the given parameters.
// A Config object represents the configuration for logging to a remote server.
// Config object is the main point interactions between user and bark client library.
//
// The url parameter is the base URL of the remote bark server where the logs will be sent.
// It must be a valid URL string and must end in `/`
//
// The defaultLogLvl parameter is the log level for logging. It must be one of the constants
// defined in the constants package, such as INFO, WARN, ERROR, etc. If an invalid value
// is given, the function will print a warning message and use INFO as the default level.
//
// The svcName parameter is the name of the service that is logging. It must be a non-empty
// string. If an empty string is given, the function will print a warning message and use
// constants.DefaultLogServiceName as the default value.
//
// The svcInstName parameter is the name of the service instance that is logging. It must be a non-empty
// string. If an empty string is given, the function will print a warning message and use
// appRuntime.SessionName as the default value.
//
// The enableSlog parameter is a boolean flag that indicates whether to enable slog logging
// to standard output. If true, the function will create and assign a new slog.Logger object
// to the Config object. If false, the Config object will have a nil Slogger field.
//
// The enableBulkSend parameter is a boolean flag that indicates whether to enable bulk sending
// of logs to the remote server. If true, the function will start a goroutine that periodically
// sends all the buffered logs to the server. If false, the logs will be sent individually as
// they are generated.
func NewClient(url, defaultLogLvl, svcName, svcInstName string, enableSlog bool, enableBulkSend bool) *Config {
	if !isValid(defaultLogLvl) {
		fmt.Printf("L#1LPYG2 - %v is not an acceptable log level. %v will be used as the default log level", defaultLogLvl, constants.DefaultLogLevel)
		defaultLogLvl = constants.DefaultLogLevel
	}

	if strings.TrimSpace(svcName) == "" {
		svcName = constants.DefaultLogServiceName
		fmt.Printf("L#1L3WBF - Blank service name supplied. Using %v as Service Name", svcName)
	}

	if strings.TrimSpace(svcInstName) == "" {
		svcInstName = appRuntime.SessionName
		fmt.Printf("L#1L3WBF - Blank instance name supplied. Using %v as Service Instance Name", svcInstName)
	}

	//Wg.Add(1)

	if enableBulkSend {
		go keepSendingLogs(url)
	}

	var slogger *slog.Logger

	if enableSlog {
		slogger = newSlogger(os.Stdout)
	} else {
		slogger = nil
	}

	return &Config{
		serverMode:          constants.ClientServerUsageModeRemoteServer,
		BaseUrl:             url,
		ErrorLevel:          defaultLogLvl,
		ServiceName:         svcName,
		ServiceInstanceName: svcInstName,
		Slogger:             slogger,
		BulkSend:            enableBulkSend,
	}
}

// NewClientWithServer returns a client config which performs the job of the server as well
// It differs from NewClient in two main ways: it does not have the option to do bulk inserts (they are not needed)
// and it accepts the database URL instead of server URL.
//
// The url parameter is the database URL where the logs will be stored.
// It must be a valid postgresql protocol string.
//
// The defaultLogLvl parameter is the log level for logging. It must be one of the constants
// defined in the constants package, such as INFO, WARN, ERROR, etc. If an invalid value
// is given, the function will print a warning message and use INFO as the default level.
//
// The svcName parameter is the name of the service that is logging. It must be a non-empty
// string. If an empty string is given, the function will print a warning message and use
// constants.DefaultLogServiceName as the default value.
//
// The svcInstName parameter is the name of the service instance that is logging. It must be a non-empty
// string. If an empty string is given, the function will print a warning message and use
// appRuntime.SessionName as the default value.
//
// The enableSlog parameter is a boolean flag that indicates whether to enable slog logging
// to standard output. If true, the function will create and assign a new slog.Logger object
// to the Config object. If false, the Config object will have a nil Slogger field.
func NewClientWithServer(dbUrl, defaultLogLvl, svcName, svcInstName string, enableSlog bool) *Config {
	if !isValid(defaultLogLvl) {
		fmt.Printf("L#1M1XXN - %v is not an acceptable log level. %v will be used as the default log level", defaultLogLvl, constants.DefaultLogLevel)
		defaultLogLvl = constants.DefaultLogLevel
	}

	if strings.TrimSpace(svcName) == "" {
		svcName = constants.DefaultLogServiceName
		fmt.Printf("L#1M1XY9 - Blank service name supplied. Using %v as Service Name", svcName)
	}

	if strings.TrimSpace(svcInstName) == "" {
		svcInstName = appRuntime.SessionName
		fmt.Printf("L#1M1XZH - Blank instance name supplied. Using %v as Service Instance Name", svcInstName)
	}

	var slogger *slog.Logger

	if enableSlog {
		slogger = newSlogger(os.Stdout)
	} else {
		slogger = nil
	}

	// Connect to the database
	err := utils.ParsePostgresUrl(dbUrl)
	if err != nil {
		panic(err.Error())
	}

	err = resources.InitDb(dbUrl)
	if err != nil {
		panic("E#1M2UFR - " + err.Error())
	}

	bld := models.NewBarkLogDao()
	err = bld.InsertServerStartedLog()
	if err != nil {
		panic("P#1LQ2YQ - Bark server start failed: " + err.Error())
	}

	// Start the server side
	go dbLogWriter.KeepSavingLogs()

	return &Config{
		serverMode:          constants.ClientServerUsageModeEmbedded,
		BaseUrl:             constants.DisabledServerUrl,
		ErrorLevel:          defaultLogLvl,
		ServiceName:         svcName,
		ServiceInstanceName: svcInstName,
		Slogger:             slogger,
		BulkSend:            false,
	}
}

// NewSloggerClientJson returns a Config object that will print logs to stdout in JSON format.
// *Config returned by NewSloggerClientJson behaves the same as NewSloggerClient, but it prints the logs in JSON.
// NewClientWithJSONSlogger accepts the same parameters as NewSloggerClient.
func NewSloggerClientJson(defaultLogLvl string) *Config {
	client := NewSloggerClient(defaultLogLvl)
	client.SetSlogHandler(slog.NewJSONHandler(os.Stdout, SlogHandlerOptions()))
	return client
}

// SetCustomOut allows users to set output to custom writer instead of the default standard output
func (c *Config) SetCustomOut(out io.Writer) {
	c.Slogger = newSlogger(out)
}

// SetSlogHandler allows users to specify their own slog handler
func (c *Config) SetSlogHandler(handler slog.Handler) {
	c.Slogger = newSlogWithCustomHandler(handler)
}

// WaitAndEnd will wait for all logs to be sent to server.
// This is an optional blocking call if the there are unsent logs.
func (c *Config) WaitAndEnd() {
	switch c.serverMode {
	case constants.ClientServerUsageModeRemoteServer:
		Wg.Wait()
	case constants.ClientServerUsageModeEmbedded:
		// Server is embedded in the client. Wait for the server side saver wait group to finish
		resources.ServerDbSaverWg.Wait()
	case constants.ClientServerUsageModeDisabled:
		// nothing to do
		return
	default:
		// This is not supposed to happen!
		panic("P#1M2YQ4 - Invalid server mode for client")
	}
}
