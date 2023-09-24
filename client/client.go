package client

import (
	"fmt"
	"github.com/techrail/bark/models"
)

type Config struct {
	BaseUrl     string
	ErrorLevel  string
	ServiceName string
	SessionName string
}

func (c *Config) Error(message string) {
	c.log(message, "ERROR")
}

func (c *Config) Info(message string) {
	c.log(message, "INFO")
}

func (c *Config) Warn(message string) {
	c.log(message, "WARN")
}

func (c *Config) Debug(message string) {
	c.log(message, "DEBUG")
}

func (c *Config) Errorf(message string, format ...any) {
	c.log(fmt.Sprintf(message, format...), "ERROR")
}

func (c *Config) Infof(message string, format ...any) {
	c.log(fmt.Sprintf(message, format...), "INFO")
}

func (c *Config) Warnf(message string, format ...any) {
	c.log(fmt.Sprintf(message, format...), "WARN")
}

func (c *Config) Debugf(message string, format ...any) {
	c.log(fmt.Sprintf(message, format...), "DEBUG")
}

func (c *Config) log(message, logLevel string) {
	// Todo: We have to parse the error message
	log := models.BarkLog{
		Message:     message,
		LogLevel:    logLevel,
		SessionName: c.SessionName,
		ServiceName: c.ServiceName,
	}

	log.Code = getCode(&log)

	go func() {
		_, err := PostLog(c.BaseUrl+"/insertSingle", log)
		if err.Severity == 1 {
			fmt.Println(err.Error())
			return
		}
	}()

	fmt.Printf("%s:\t %s -- %s\n", logLevel, c.SessionName, message)
	// Todo: Add uber zap to avoid printing with PrintF (We don't want to handle log printing)
}

func getCode(log *models.BarkLog) string {
	// Todo: Generate an error code like E#ERRCODE where "E" indicate loglevel
	return "00000"
}

func NewClient(url, errLevel, svcName, sessName string) *Config {
	return &Config{
		BaseUrl:     url,
		ErrorLevel:  errLevel,
		ServiceName: svcName,
		SessionName: sessName,
	}
}
