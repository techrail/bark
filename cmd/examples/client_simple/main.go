package main

import (
	"fmt"
	"github.com/techrail/bark/client"
	"os"
	"time"
)

func main() {
	err := os.Setenv("DEBUG", "true")
	if err != nil {
		fmt.Println("This")
		return
	}
	log := client.NewClient("http://127.0.0.1:8080/", "INFO", "brktest", "local session", true, true)

	log.Panic("Panic message")
	log.Alert("Alert message", true)
	log.Error("Error message")
	log.Warn("Warn message")
	log.Notice("Notice message")
	log.Info("Info message")
	log.Debug("Debug message")
	log.Println("Println message")

	log.Panicf("Panic message: %v", "Something")
	log.Alertf("Alert message: %v", true, "Something")
	log.Errorf("Error message: %v", "Something")
	log.Warnf("Warn message: %v", "Something")
	log.Noticef("Notice message: %v", "Something")
	log.Infof("Info message: %v", "Something")
	log.Debugf("Debug message: %v", "Something")

	log.Panic("E#1LPW24 - Panic message")
	log.Alert("E#1LPW25 - Alert message", false)
	log.Error("E#1LPW26 - Error message")
	log.Warn("E#1LPW27 - Warn message")
	log.Notice("E#1LPW28 - Notice message")
	log.Info("E#1LPW29 - Info message")
	log.Debug("E#1LPW30 - Debug message")
	log.Println("E#1LPW30 - Println message")
	log.Printf("E#1LPW30 - Printf message")

	log.Println("P#1LPW5I - Panic message")
	log.Println("A#1LPW5O - Alert message")
	log.Println("E#1LPW5Q - Error message")
	log.Println("W#1LPW5V - Warn message")
	log.Println("N#1LPW5Y - Notice message")
	log.Println("I#1LPW61 - Info message")
	log.Println("D#1LPW65 - Debug message")

	log.Default("P#1LTOTD - Panic message")
	log.Default("A#1LTOTK - Alert message")
	log.Default("E#1LTOTN - Error message")
	log.Default("W#1LTOTQ - Warn message")
	log.Default("N#1LTOTT - Notice message")
	log.Default("I#1LTOTW - Info message")
	log.Default("D#1LTOU0 - Debug message")
	log.Default("1LTOU4 - Default message")

	log.Println("D # 1LPW68 - Debug message")
	log.Println("")
	log.Println(" ")

	log.Println("1LPWC7 - Default message")

	time.Sleep(1 * time.Second)
	// --------------------------------
	//time.Sleep(1 * time.Second)
}
