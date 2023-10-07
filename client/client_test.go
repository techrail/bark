package client

import (
	"fmt"
	"github.com/techrail/bark/models"
	"testing"

	"github.com/techrail/bark/constants"
)

var client *Config

func init() {
	client = NewClient("http://localhost/bark", "info", "testSvc", "testSess")
}

func TestConfigParseMessage(t *testing.T) {
	l := client.parseMessage("D#LMID1 - Log message1")
	if l.LogLevel != constants.Debug {
		t.Errorf("Expected Debug level")
	} else {
		fmt.Println("1...OK")
	}

	if l.Code != "LMID1" {
		t.Errorf("Expected LMID1")
	} else {
		fmt.Println("2...OK")
	}

	if l.Message != "Log message1" {
		t.Errorf("Should be `Log message1`. Found `%v`", l.Message)
	} else {
		fmt.Println("3...OK")
	}

	// ------------------------------------------------------------
	l = client.parseMessage("I#LMID2 - Log message2")
	if l.LogLevel != constants.Info {
		t.Errorf("Expected Info level")
	} else {
		fmt.Println("4...OK")
	}

	if l.Code != "LMID2" {
		t.Errorf("Expected LMID2")
	} else {
		fmt.Println("5...OK")
	}

	if l.Message != "Log message2" {
		t.Errorf("Should be `Log message2`. Found `%v`", l.Message)
	} else {
		fmt.Println("6...OK")
	}

	// --------------------------------------------------
	l = client.parseMessage("L# - Log message 3")
	if l.LogLevel != constants.DefaultLogLevel {
		t.Errorf("Should have been an error!")
	} else {
		fmt.Println("7...OK")
	}

	if l.Code != constants.DefaultLogCode {
		t.Errorf("Not the right code!")
	} else {
		fmt.Println("8...OK")
	}

	if l.Message != "L# - Log message 3" {
		t.Errorf("Should be `L# - Log message 3`!")
	} else {
		fmt.Println("9...OK")
	}

	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
	// --------------------------------------------------
}

func TestPrint(t *testing.T) {
	m := "D#LMID1 - Log message1"
	l := client.parseMessage(m)
	printParsedInfo(m, l)

	m = "I#LMID2 - Log message2"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "N#LMID3 - Log message3"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "W#LMID4 - Log message4"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "E#LMID5 - Log message5"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "A#LMID6 - Log message6"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "P#LMID7 - Log message7"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "LMID8 - Log message8"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "Log message9"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "- Log message10"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = " # - Log message11"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "X# - Log message12"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "XX# - Log message13"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "XX#LMID14 - Log message14"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "XX#LMIDINTHISCASEISVERYVERYLONG - Log message15"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "#LMIDINTHISCASEISVERYVERYLONG - Log message16"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "#LMID17 - Log message17"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "D#LMID18 - Log message18"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "LMIDINTHISCASEISVERYVERYLONG - Log message19"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "E#LMID20 - Log message20 - with - more - dashes"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = "XX#LMID21 - Log message21"
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = ""
	l = client.parseMessage(m)
	printParsedInfo(m, l)

	m = " "
	l = client.parseMessage(m)
	printParsedInfo(m, l)
}

func printParsedInfo(msg string, l models.BarkLog) {
	fmt.Printf("| `%v` | `%v` | `%v` | `%v` |\n", msg, l.LogLevel, l.Code, l.Message)
}
