package main

import (
	"fmt"
	"github.com/techrail/bark/client"
	"log/slog"
	"os"
)

func main() {
	log := client.NewClient("http://127.0.0.1:8080/", "INFO", "BarkClientFileTest", "TestClientSession", true, false)

	file, err := os.Create("random.txt")
	if err != nil {
		fmt.Println("E#1M5WRN - Error when creating new file: ", err)
		return
	}

	log.SetSlogHandler(slog.NewJSONHandler(file, client.SlogHandlerOptions()))

	log.Info("Some Message that'll be sent to random.txt file")
	log.WaitAndEnd()
}
