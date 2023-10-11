package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/techrail/bark/models"
	"net/http"
	"net/url"
)

type webhook struct {
	authentication Authentication
	url            url.URL
}

func (w *webhook) StartPublishing() chan<- models.BarkLog {
	ret := make(chan models.BarkLog)
	go func(ch chan models.BarkLog) {
		for {
			msg := <-ch
			body := bytes.NewBuffer([]byte{})
			json.NewEncoder(body).Encode(msg)
			req, err := http.NewRequest(http.MethodPost, w.url.String(), body)
			if err != nil {
				//todo: this
			}
			w.authentication.AddAuth(req)
			c := &http.Client{}
			_, err = c.Do(req)
			if err != nil {
				//todo: this
			}

		}
	}(ret)
	return ret
}

type Authentication interface {
	AddAuth(*http.Request)
}

// todo: add basic http auth
type HttpBasicAuth struct {
	username string
	password string
}

func NewHttpBasicAuth(username string, password string) Authentication {
	return HttpBasicAuth{
		username: username,
		password: password,
	}
}

func (b HttpBasicAuth) AddAuth(r *http.Request) {
	r.SetBasicAuth(b.username, b.password)
}

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
