package connectors

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/microlib/simple"
)

// Mock all connections
type MockConnectors struct {
	Logger *simple.Logger
	Http   *http.Client
	Type   string
}

func (c *MockConnectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Do(req *http.Request) (*http.Response, error) {
	if c.Type == "error" {
		return nil, errors.New("Forced error")
	}
	return c.Http.Do(req)
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewTestConnector(data string, code int, con string, logger *simple.Logger) Client {

	// we first load the json payload to simulate a call
	// for now just ignore failures.
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(data)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conns := &MockConnectors{Http: httpclient, Logger: logger, Type: con}
	return conns
}
