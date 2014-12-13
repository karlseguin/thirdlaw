package checks

import (
	"gopkg.in/karlseguin/typed.v1"
	"github.com/karlseguin/beats/core"
	"strings"
	"fmt"
)

type Base struct {
	name string
	runner core.Runner
}

func (c *Base) Name() string {
	return c.name
}

func (c *Base) Run() *core.Result {
	return c.runner.Run()
}

func New(t typed.Typed) core.Check {
	switch strings.ToLower(t.String("type")) {
	case "http":
		return build(t, NewHttp(t))
	default:
		b, _ := t.ToBytes("")
		panic(fmt.Errorf("missing type for  %v", string(b)))
	}
}

func build(t typed.Typed, runner core.Runner) core.Check {
	name, ok := t.StringIf("name")
	if ok == false {
		b, _ := t.ToBytes("")
		panic(fmt.Errorf("missing name %v", string(b)))
	}
	c := &Base{
		name: name,
		runner: runner,
	}
	return c
}
