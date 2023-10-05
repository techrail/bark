package client

import (
	"fmt"
	"github.com/techrail/bark/client/controllers"
	"github.com/techrail/bark/client/services/clientLogSender"
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
	// AlertWebhook webhook
}

func (c *Config) parseMessage(msg string) models.BarkLog {
	l := models.BarkLog{
		ServiceName: c.ServiceName,
		SessionName: c.SessionName,
	}
	// Look for `-` in the message
	pos := strings.Index(msg, "-")
	if pos < 1 {
		// There is no `-` in the message.
		l.Message = msg
		l.Code = constants.DefaultLogCode
		return l
	}
	// separate the message and meta info
	l.Message = msg[pos:]
	meta := msg[:pos]

	// Separate the code and level
	metas := strings.Split(meta, "#")
	if len(metas) != 2 {
		// Improperly formatted message
		l.Message = msg
	}

	logLvl := strings.TrimSpace(metas[0])
	logCode := strings.TrimSpace(metas[1])

	if len(logLvl) != 1 {
		l.LogLevel = constants.Info
	} else {
		l.LogLevel = getLogLevelFromCharacter(metas[0])
	}

	if len(logCode) < 1 || len(logCode) > 16 {
		l.Code = constants.DefaultLogMessage
	} else {
		l.Code = logCode
	}

	return l
}

func getLogLevelFromCharacter(s string) string {
	switch strings.ToUpper(s) {
	case "P":
		return constants.Error
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
	c.sendLogToServer(message, constants.Panic)
}
func (c *Config) Alert(message string) {
	// Todo: handle the alert webhook call here
	c.sendLogToServer(message, constants.Alert)
}
func (c *Config) Error(message string) {
	c.sendLogToServer(message, constants.Error)
}
func (c *Config) Warn(message string) {
	c.sendLogToServer(message, constants.Warning)
}
func (c *Config) Notice(message string) {
	c.sendLogToServer(message, constants.Notice)
}
func (c *Config) Info(message string) {
	c.sendLogToServer(message, constants.Info)
}
func (c *Config) Debug(message string) {
	c.sendLogToServer(message, constants.Debug)
}
func (c *Config) Println(message string) {
	c.sendLogToServer(message+"\n", constants.Info)
}

func (c *Config) ParseAndSend(message string) {
	controllers.SendSingleToClientChannel(c.parseMessage(message))
}

func (c *Config) Panicf(message string, format ...any) {
	c.sendLogToServer(fmt.Sprintf(message, format...), constants.Panic)
}
func (c *Config) Alertf(message string, format ...any) {
	c.sendLogToServer(fmt.Sprintf(message, format...), constants.Alert)
}
func (c *Config) Errorf(message string, format ...any) {
	c.sendLogToServer(fmt.Sprintf(message, format...), constants.Error)
}
func (c *Config) Warnf(message string, format ...any) {
	c.sendLogToServer(fmt.Sprintf(message, format...), constants.Warning)
}
func (c *Config) Noticef(message string, format ...any) {
	c.sendLogToServer(fmt.Sprintf(message, format...), constants.Notice)
}
func (c *Config) Infof(message string, format ...any) {
	c.sendLogToServer(fmt.Sprintf(message, format...), constants.Info)
}
func (c *Config) Debugf(message string, format ...any) {
	c.sendLogToServer(fmt.Sprintf(message, format...), constants.Debug)
}

// func (c *Config) SetAlertWebhook(f webhook) {
// 	c.AlertWebhook = f
// }

func (c *Config) sendLogToServer(message, logLevel string) {
	log := models.BarkLog{
		Message:     message,
		LogLevel:    logLevel,
		SessionName: c.SessionName,
		ServiceName: c.ServiceName,
	}

	log.Code = getCode(&log)

	controllers.SendSingleToClientChannel(log)

	fmt.Printf("%s:\t %s -- %s\n", logLevel, c.SessionName, message)
	// Todo: Add uber zap to avoid printing with PrintF (We don't want to handle sendLogToServer printing)
}

func getCode(log *models.BarkLog) string {
	// Todo: Generate an error code like E#ERRCODE where "E" indicate loglevel
	return "00000"
}

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
	}
}
