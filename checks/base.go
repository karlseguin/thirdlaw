package checks

import (
	"fmt"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"strings"
	"log"
	"time"
)

type Base struct {
	name   string
	onFail []core.Action
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
	if res.Ok == false {
		for _, fail := range c.onFail {
			if err := fail.Run(); err != nil {
				log.Println(err)
			}
		}
	}
	return res
}

func New(actions map[string]core.Action, t typed.Typed) core.Check {
	switch strings.ToLower(t.String("type")) {
	case "http":
		return build(actions, t, NewHttp(t))
	default:
		panic(fmt.Errorf("unknown type %v", string(t.MustBytes(""))))
	}
}

func build(actions map[string]core.Action, t typed.Typed, runner core.Runner) core.Check {
	name, ok := t.StringIf("name")
	if ok == false {
		panic(fmt.Errorf("missing name %v", string(t.MustBytes(""))))
	}
	c := &Base{
		name:   name,
		runner: runner,
	}
	if failNames := t.Strings("onFail"); len(failNames) > 0 {
		c.onFail = make([]core.Action, len(failNames))
		for i, n := range failNames {
			action, ok := actions[n]
			if ok == false {
				panic(fmt.Errorf("unknown action %q for check %q", n, name))
			}
			c.onFail[i] = action
		}
	}
	return c
}
