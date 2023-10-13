package channels

import (
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/models"
)

var LogChannel chan models.BarkLog

func init() {
	LogChannel = make(chan models.BarkLog, constants.ServerLogInsertionChannelCapacity)
}
