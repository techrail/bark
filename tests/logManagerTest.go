package tests

import (
	"github.com/techrail/bark/logManager"
)

func TestInfoLogger() {
	logger := logManager.GetLogger("", "")

	logger.Info("test", "00000")
	logger.Info("test2", "00000", map[string]string{"caramel": "chocolate"})
}
