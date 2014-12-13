package beats

import (
	"gopkg.in/karlseguin/typed.v1"
	"time"
)

type CheckHttp struct {
	name    string
	host    string
	path    string
	timeout time.Duration
}

func (c *CheckHttp) Name() string {
	return c.name
}

func (c *CheckHttp) Run() *Result {
	return OKResult
}

func NewHttpCheck(name string, t typed.Typed) Check {
	return &CheckHttp{
		name:    name,
		host:    t.StringOr("host", "127.0.0.1"),
		path:    t.StringOr("path", "/"),
		timeout: time.Millisecond * time.Duration(t.IntOr("timeout", 5000)),
	}
}
