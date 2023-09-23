package tests

import (
	"github.com/techrail/bark/client"
	"github.com/techrail/bark/models"
	"testing"
)

func Test_requester(t *testing.T) {
	var log models.BarkLog
	log.Message = "Some random message"
	response, err := client.PostLog("http://localhost:8080/insertSingle", log)
	if err.Severity == 1 {
		t.Error("Error: " + err.Error())
	}
	t.Log(response + " response of inserting a log")
}
