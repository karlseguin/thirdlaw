package checks

import (
	"bytes"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"io/ioutil"
	"net/http"
	"time"
)

type Http struct {
	address  string
	contains []byte
	client   *http.Client
}

func (c *Http) Run() *core.Result {
	res, err := c.client.Get(c.address)
	if err != nil {
		return core.Error(err)
	}
	defer res.Body.Close()
	if res.StatusCode < 300 && c.contains == nil {
		return core.Success()
	}
	data, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode < 300 {
		if bytes.Contains(data, c.contains) {
			return core.Success()
		}
		return core.Failuref("response did not contain expected content")
	}
	return core.Failuref("%d: %s", res.StatusCode, string(data))
}

func NewHttp(t typed.Typed) *Http {
	var contains []byte
	if c, ok := t["contains"]; ok {
		contains = []byte(c.(string))
	}
	return &Http{
		contains: contains,
		address:  t.StringOr("address", "http://127.0.0.1/"),
		client: &http.Client{
			Timeout: time.Second * time.Duration(t.IntOr("timeout", 5)),
		},
	}
}
