package xhttputil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Transport struct {
	handlerBody func(body string) string
}

func NewTransport(handlerBody func(body string) string) *Transport {
	return &Transport{
		handlerBody: handlerBody,
	}
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(t.handlerBody(string(b)))))
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}
