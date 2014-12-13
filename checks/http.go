package checks

import (
	"github.com/karlseguin/beats/core"
	"gopkg.in/karlseguin/typed.v1"
	"time"
)

type Http struct {
	host    string
	path    string
	timeout time.Duration
}

func (c *Http) Run() *core.Result {
	return core.Success()
}

func NewHttp(t typed.Typed) *Http {
	return &Http{
		host:    t.StringOr("host", "127.0.0.1"),
		path:    t.StringOr("path", "/"),
		timeout: time.Millisecond * time.Duration(t.IntOr("timeout", 5000)),
	}
}
