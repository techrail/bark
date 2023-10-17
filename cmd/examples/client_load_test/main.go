package main

import (
	"fmt"
	"github.com/techrail/bark/client"
	"os"
)

func main() {
	err := os.Setenv("DEBUG", "true")
	if err != nil {
		fmt.Println("This")
		return
	}
	log := client.NewClient("http://127.0.0.1:8080/", "INFO", "brktest", "load test session", false, true)

	for i := 0; i < 5_000_000; i++ {
		log.Printf("1LTOU4 - Default message - %v", i)
	}

	log.WaitAndEnd()
}
