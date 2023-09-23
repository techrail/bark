package tests

import (
	"github.com/techrail/bark/logManager"
	"testing"
)

func TestGetLogger(t *testing.T) {
	logger := logManager.GetLogger("http://localhost:8080/", "JustRandomName")
	logger.Info("Info")
}
