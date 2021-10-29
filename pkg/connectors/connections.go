package connectors

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	Logger *simple.Logger
	Http   *http.Client
}

func NewClientConnections(logger *simple.Logger) Client {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	return &Connectors{Http: httpClient, Logger: logger}
}

func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Do(req *http.Request) (*http.Response, error) {
	return c.Http.Do(req)
}
