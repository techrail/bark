package client

import (
	"fmt"
	"github.com/techrail/bark/constants"
	"testing"
)

var client *Config

func init() {
	client = NewClient("http://localhost/bark", "info", "testSvc", "testSess")
}

func TestConfigParseMessage(t *testing.T) {
	l := client.parseMessage("E#1LFV5T - Log message")
	if l.LogLevel != constants.Error {
		t.Errorf("1LFVD5 - Should have been an error!")
	} else {
		fmt.Println("...OK")
	}

	if l.Code != "1LFV5T" {
		t.Errorf("1LFVOH - Not the right code!")
	} else {
		fmt.Println("...OK")
	}

	if l.Message != "Log message" {
		t.Errorf("E#1LFVSV - Should be `Log message`!")
	} else {
		fmt.Println("...OK")
	}

	// ------------------------------------------------------------
	l := client.parseMessage("E#1LFV5T - Log message")
	if l.LogLevel != constants.Error {
		t.Errorf("1LFVD5 - Should have been an error!")
	} else {
		fmt.Println("...OK")
	}

	if l.Code != "1LFV5T" {
		t.Errorf("1LFVOH - Not the right code!")
	} else {
		fmt.Println("...OK")
	}

	if l.Message != "Log message" {
		t.Errorf("E#1LFVSV - Should be `Log message`!")
	} else {
		fmt.Println("...OK")
	}

}
