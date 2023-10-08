package alerts

import (
	"fmt"
	"net/http"
)

type WebhookMessage struct {
	ServiceName string
	SessionName string
	Message     string
}

type Webhook interface {
	StartPublishing(url string) chan WebhookMessage
}

type Authentication interface {
	AddAuth(*http.Request)
}

//todo: add basic http auth

type BearerAuth struct {
	bearerToken string
	headerName  string
	prefix      string
}

func NewBearerAuth(bearerToken string, headerName string, prefix string) Authentication {
	return BearerAuth{
		bearerToken: bearerToken,
		headerName:  headerName,
		prefix:      prefix,
	}
}

func (b BearerAuth) AddAuth(r *http.Request) {
	if b.prefix == "" {
		r.Header.Add(b.headerName, b.bearerToken)
	} else {
		r.Header.Add(b.headerName, fmt.Sprintf("%s %s", b.prefix, b.bearerToken))
	}
}
