package client

import (
	"fmt"
	"github.com/techrail/bark/models"
	"testing"

	"github.com/techrail/bark/constants"
)

var client *Config

func init() {
	client = NewClient("http://localhost/bark", "info", "testSvc", "testSess", true, false)
}

func TestConfigParseMessage1(t *testing.T) {
	l := client.parseMessage("D#LMID1 - Log message1")
	if l.LogLevel != constants.Debug || l.Code != "LMID1" || l.Message != "Log message1" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage2(t *testing.T) {
	l := client.parseMessage("I#LMID2 - Log message2")
	if l.LogLevel != constants.Info || l.Code != "LMID2" || l.Message != "Log message2" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage3(t *testing.T) {
	l := client.parseMessage("N#LMID3 - Log message3")
	if l.LogLevel != constants.Notice || l.Code != "LMID3" || l.Message != "Log message3" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage4(t *testing.T) {
	l := client.parseMessage("W#LMID4 - Log message4")
	if l.LogLevel != constants.Warning || l.Code != "LMID4" || l.Message != "Log message4" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage5(t *testing.T) {
	l := client.parseMessage("E#LMID5 - Log message5")
	if l.LogLevel != constants.Error || l.Code != "LMID5" || l.Message != "Log message5" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage6(t *testing.T) {
	l := client.parseMessage("A#LMID6 - Log message6")
	if l.LogLevel != constants.Alert || l.Code != "LMID6" || l.Message != "Log message6" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage7(t *testing.T) {
	l := client.parseMessage("P#LMID7 - Log message7")
	if l.LogLevel != constants.Panic || l.Code != "LMID7" || l.Message != "Log message7" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage8(t *testing.T) {
	l := client.parseMessage("LMID8 - Log message8")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != "LMID8" || l.Message != "Log message8" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage9(t *testing.T) {
	l := client.parseMessage("Log message9")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "Log message9" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage10(t *testing.T) {
	l := client.parseMessage("- Log message10")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "- Log message10" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage11(t *testing.T) {
	l := client.parseMessage(" # - Log message11")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != " # - Log message11" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage12(t *testing.T) {
	l := client.parseMessage("X# - Log message12")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "X# - Log message12" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage13(t *testing.T) {
	l := client.parseMessage("XX# - Log message13")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "XX# - Log message13" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage14(t *testing.T) {
	l := client.parseMessage("XX#LMID14 - Log message14")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "XX#LMID14 - Log message14" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage15(t *testing.T) {
	l := client.parseMessage("XX#LMIDINTHISCASEISVERYVERYLONG - Log message15")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "XX#LMIDINTHISCASEISVERYVERYLONG - Log message15" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage16(t *testing.T) {
	l := client.parseMessage("#LMIDINTHISCASEISVERYVERYLONG - Log message16")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "#LMIDINTHISCASEISVERYVERYLONG - Log message16" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage17(t *testing.T) {
	l := client.parseMessage("#LMID17 - Log message17")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "#LMID17 - Log message17" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage18(t *testing.T) {
	l := client.parseMessage("D#LMID18 - Log message18")
	if l.LogLevel != constants.Debug || l.Code != "LMID18" || l.Message != "Log message18" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage19(t *testing.T) {
	l := client.parseMessage("LMIDINTHISCASEISVERYVERYLONG - Log message19")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "LMIDINTHISCASEISVERYVERYLONG - Log message19" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage20(t *testing.T) {
	l := client.parseMessage("E#LMID20 - Log message20 - with - more - dashes")
	if l.LogLevel != constants.Error || l.Code != "LMID20" || l.Message != "Log message20 - with - more - dashes" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage21(t *testing.T) {
	l := client.parseMessage("XX#LMID21 - Log message21")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "XX#LMID21 - Log message21" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage22(t *testing.T) {
	l := client.parseMessage("")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != "" {
		t.Errorf("unexpected values")
	}
}

func TestConfigParseMessage23(t *testing.T) {
	l := client.parseMessage(" ")
	if l.LogLevel != constants.DefaultLogLevel || l.Code != constants.DefaultLogCode || l.Message != " " {
		t.Errorf("unexpected values")
	}
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
