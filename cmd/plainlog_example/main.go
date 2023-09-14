package main

import (
	`log`
	`time`
)

func main() {
	log.Println("Server started at " + time.Now().Format(time.RFC3339))
}
