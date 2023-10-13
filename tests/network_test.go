package tests

import (
	"github.com/techrail/bark/client"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/models"
	"log/slog"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	logClient := client.NewClient("http://localhost:8080/", constants.Info, "TestService", "localRun", false, false)
	sendDummyLogsThroughAllMethods(logClient)
}

func TestClientWithCustomOut(t *testing.T) {
	logClient := client.NewClient("http://localhost:8080/", constants.Info, "TestService", "localRun", false, false)

	file, _ := os.Create("../tmp/random.txt")

	logClient.WithCustomOut(file)

	sendDummyLogsThroughAllMethods(logClient)
}

func TestClientWithCustomSlogHandler(t *testing.T) {
	logClient := client.NewClient("http://localhost:8080/", constants.Info, "TestService", "localRun", false, false)

	file, _ := os.Create("../tmp/random.txt")

	logClient.WithSlogHandler(slog.NewJSONHandler(file, client.Options()))

	sendDummyLogsThroughAllMethods(logClient)
}

func TestPostLogArray(t *testing.T) {
	var logs []models.BarkLog
	logs = make([]models.BarkLog, 3)
	logs[0] = models.BarkLog{Message: "someMessage"}
	logs[1] = models.BarkLog{Message: "someMessage"}
	logs[2] = models.BarkLog{Message: "someMessage"}
	_, _ = client.PostLogArray("http://localhost:8080/"+constants.BatchInsertUrl, logs)
}

func sendDummyLogsThroughAllMethods(logClient *client.Config) {
	// Print without formatter
	logClient.Error("Error Message!")
	logClient.Info("Info Message!")
	logClient.Debug("Debug Message!")
	logClient.Panic("Panic Message!")
	logClient.Alert("Alert Message!", true)
	logClient.Notice("Notice Message!")

	// Print with formatter
	logClient.Errorf("%s Message!", "Error")
	logClient.Infof("%s Message!", "Info")
	logClient.Debugf("%s Message!", "Debug")
	logClient.Panicf("%s Message!", "Panic")
	logClient.Alertf("%s Message!", true, "Alert")
	logClient.Noticef("%s Message!", "Notice")
}
