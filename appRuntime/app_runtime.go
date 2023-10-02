package appRuntime

import (
	`os`
	"sync/atomic"

	`github.com/techrail/bark/utils`
)

var ShutdownRequested atomic.Bool
var SessionName string
var err error

func init() {
	ShutdownRequested.Store(false)
	SessionName, err = os.Hostname()
	if err != nil {
		SessionName = utils.GetRandomAlphaString(32)
	}
}
