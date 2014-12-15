package checks

import (
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"io/ioutil"
	"net/http"
	"time"
)

type Http struct {
	address string
	client  *http.Client
}

func (c *Http) Run() *core.Result {
	res, err := c.client.Get(c.address)
	if err != nil {
		return core.Error(err)
	}
	defer res.Body.Close()
	if res.StatusCode < 300 {
		return core.Success()
	}
	data, _ := ioutil.ReadAll(res.Body)
	return core.Failuref("%d: %s", res.StatusCode, string(data))
}

func NewHttp(t typed.Typed) *Http {
	return &Http{
		address: t.StringOr("address", "http://127.0.0.1:3000/"),
		client: &http.Client{
			Timeout: time.Millisecond * time.Duration(t.IntOr("timeout", 5000)),
		},
	}
}
