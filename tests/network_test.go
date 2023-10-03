package tests

import (
	"github.com/techrail/bark/client"
	"github.com/techrail/bark/models"
	"testing"
)

func Test_requester(t *testing.T) {
	logClient := client.NewClient("http://localhost:8080/", "INFO", "ServicName", "localRun")

	// Print with formatter

	// Print with formatter
	logClient.Error("Anime: Naruto")
	logClient.Info("Anime: One Piece")
	//logClient.Debug("Anime: Bleach")
	//logClient.Warn("Anime: AOT")
	//
	//// Print without formatter
	//logClient.Errorf("Anime: %s", "Full Metal Alchemist")
	//logClient.Infof("Anime: %s", "Tokyo Ghoul")
	//logClient.Warnf("Anime: %s", "")
	//logClient.Debugf("I want to print something! %s", "weirdString")
	//
	//// Multiple Logs
	var logs []models.BarkLog
	logs = make([]models.BarkLog, 3)
	logs[0] = models.BarkLog{Message: "someMessage"}
	logs[1] = models.BarkLog{Message: "someMessage"}
	logs[2] = models.BarkLog{Message: "someMessage"}
	logClient.Debug("Random error")
}
