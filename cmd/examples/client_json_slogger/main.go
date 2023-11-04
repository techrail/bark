package main

import "github.com/techrail/bark/client"

func main() {
	log := client.NewClientWithJSONSlogger("http://127.0.0.1:8080/", "INFO", "BarkClientFileTest", "TestClientSession", false)
	log.Info("1N09FW - This is an Info message!")
	log.Debug("1N09GG - This is an Debug message!")
	log.Warn("1N09H5 - This is an Warn message!")
	log.Notice("1N09HL - This is an Notice message!")
	log.Error("1N09HT - This is an Error message!")
	log.Panic("1N09I7 - This is an Panic message!")
	log.Alert("1N09IG - This is an Alert message!", false)
	log.Default("I#1N09JH - This is an Default message!")
	log.Println("D#1N09JR - This is an Print message!")
	log.WaitAndEnd()
}
