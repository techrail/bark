package appRuntime

import (
	"sync/atomic"
)

var ShutdownRequested atomic.Bool

func init() {
	ShutdownRequested.Store(false)
}
