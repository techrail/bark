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
	log := client.NewClientWithServer("postgres://vkaushal288:vkaushal288@127.0.0.1:5432/bark", "INFO", "brktest", "load test session", false)

	for i := 0; i < 500; i++ {
		log.Printf("1M2UBT - Default message - %v", i)
	}

	log.WaitAndEnd()
}
