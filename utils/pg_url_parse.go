package utils

import (
	"errors"
	nurl "net/url"
	"strings"
)

func ParsePostgresUrl(dbUrl string) error {
	if strings.TrimSpace(dbUrl) == "" {
		return errors.New("P#1LQ32D - Database URL is required")
	} else {
		u, err := nurl.Parse(dbUrl)
		if err != nil {
			return errors.New("P#1LQ36U - Database URL is not OK: " + dbUrl)
		}

		if u.Scheme != "postgres" && u.Scheme != "postgresql" {
			return errors.New("P#1LQ37D - Database URL must begin with postgres:// or postgresql:// : " + dbUrl)
		}
	}

	return nil
}
