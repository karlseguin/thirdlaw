package checks

import (
	"fmt"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"strings"
	"time"
)

type Base struct {
	name   string
	runner core.Runner
}

func (c *Base) Name() string {
	return c.name
}

func (c *Base) Run() *core.Result {
	s := time.Now()
	res := c.runner.Run()
	res.Name = c.name
	res.Milliseconds = int(time.Now().Sub(s).Nanoseconds() / 1000000)
	return res
}

func New(t typed.Typed) core.Check {
	switch strings.ToLower(t.String("type")) {
	case "http":
		return build(t, NewHttp(t))
	default:
		b, _ := t.ToBytes("")
		panic(fmt.Errorf("unknown type %v", string(b)))
	}
}

func build(t typed.Typed, runner core.Runner) core.Check {
	name, ok := t.StringIf("name")
	if ok == false {
		b, _ := t.ToBytes("")
		panic(fmt.Errorf("missing name %v", string(b)))
	}
	c := &Base{
		name:   name,
		runner: runner,
	}
	return c
}
