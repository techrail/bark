package constants

// Constants that are used mostly on the server side (but can also be used on the client side)

const (
	AppName    = "Bark"
	AppVersion = "1.0.0"
)

const (
	DefaultLogCode        = "000000"
	DefaultLogLevel       = "INFO"
	DefaultLogMessage     = "_no_msg_supplied_"
	DefaultLogServiceName = "def_svc"
	DefaultLogSessionName = "def_sess"
)

const (
	MaxLogCodelength = 16 // DB constraint
)

const (
	Panic   = "PANIC"
	Alert   = "ALERT"
	Error   = "ERROR"
	Warning = "WARN"
	Notice  = "NOTICE"
	Info    = "INFO"
	Debug   = "DEBUG"
)

const (
	SingleInsertUrl = "insertSingle"
	BatchInsertUrl  = "insertMultiple"
)
