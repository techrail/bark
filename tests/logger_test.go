package tests

import (
	"encoding/json"
	"testing"

	"github.com/techrail/bark/db"
	"github.com/techrail/bark/logger"
)

func TestGetLogger(t *testing.T) {
	serviceName := "TestService"
	bark := logger.GetLogger(serviceName)

	weirdJsonData, _ := json.Marshal(map[string]interface{}{
		"weird_key_1": "destroy the world",
		"weird_key_2": "or may wait until tomorrow",
	},
	)

	bark(0, "901101", "This logger rocks!", weirdJsonData)

	database, err := db.ConnectToDatabase()

	if err != nil {
		t.Log("Failed to connect to database for verificaiton")
		t.Fail()
	}

	logs, err := database.FetchLimitedLogs(100)

	for i := 0; i < 100; i++ {
		t.Log(logs[i].Message)
		if logs[i].Message == "This logger rocks!" {
			return
		}
	}
	t.Log("Unable to find the inserted log")
	t.Fail()
}
