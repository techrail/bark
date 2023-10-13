package appError

import (
	"fmt"
)

type AppErr struct {
	Severity int
	Code     string
	Msg      string
}

func (ae AppErr) Error() string {
	return fmt.Sprintf("E#1L3TGS - %v, %v, %v", ae.Severity, ae.Code, ae.Msg)
}

func (ae AppErr) String() string {
	return ae.Error()
}
