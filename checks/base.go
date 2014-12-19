package checks

import (
	"fmt"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"strings"
	"time"
)

type Base struct {
	name   string
	onFail []string
	runner core.Runner
	ap     core.ActionProvider
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
		for _, actionName := range c.onFail {
			action := c.ap.GetAction(actionName)
			if action == nil {
				log.Println(fmt.Sprintf("fail action %q is unknown for %q check", actionName, c.name))
			} else if err := action.Run(); err != nil {
				log.Println(err)
			}
		}
	}
	return res
}

func New(ap core.ActionProvider, t typed.Typed) core.Check {
	switch strings.ToLower(t.String("type")) {
	case "http":
		return build(ap, t, NewHttp(t))
	case "shell":
		return build(ap, t, NewShell(t))
	default:
		panic(fmt.Errorf("unknown type %v", string(t.MustBytes(""))))
	}
}

func build(ap core.ActionProvider, t typed.Typed, runner core.Runner) core.Check {
	name, ok := t.StringIf("name")
	if ok == false {
		panic(fmt.Errorf("missing name %v", string(t.MustBytes(""))))
	}
	return &Base{
		ap:     ap,
		name:   name,
		runner: runner,
		onFail: t.Strings("onFail"),
	}
}
