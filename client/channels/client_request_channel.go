package channels

import "github.com/techrail/bark/models"

const ClientChannelCapacity = 10000

var ClientChannel chan models.BarkLog

func init() {
	ClientChannel = make(chan models.BarkLog, ClientChannelCapacity)
}
