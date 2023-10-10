package client

const (
	PANIC   = "PANIC"
	ALERT   = "ALERT"
	ERROR   = "ERROR"
	WARNING = "WARN"
	NOTICE  = "NOTICE"
	INFO    = "INFO"
	DEBUG   = "DEBUG"
)

func isValid(lvl string) bool {
	switch lvl {
	case PANIC:
		fallthrough
	case ALERT:
		fallthrough
	case ERROR:
		fallthrough
	case WARNING:
		fallthrough
	case NOTICE:
		fallthrough
	case INFO:
		fallthrough
	case DEBUG:
		return true
	default:
		return false
	}
}
