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
	name      string
	recover   []string
	runner    core.Runner
	next      time.Time
	snooze    time.Duration
	ap        core.ActionProvider
	failCount int
	failures  int
}

func (c *Base) Name() string {
	return c.name
}

func (c *Base) Run() *core.Result {
	s := time.Now()
	if s.Before(c.next) {
		return nil
	}

	res := c.runner.Run()
	res.Name = c.name
	res.Milliseconds = int(time.Now().Sub(s).Nanoseconds() / 1000000)
	c.next = s.Add(c.snooze)
	if res.Ok {
		c.failures = 0
		return res
	}
	c.failures++
	res.Failures = c.failures
	if c.failures%c.failCount != 0 {
		res.Silent = true
		return res
	}
	for _, actionName := range c.recover {
		action := c.ap.GetAction(actionName)
		if action == nil {
			log.Println(fmt.Sprintf("fail action %q is unknown for %q check", actionName, c.name))
		} else if err := action.Run(); err != nil {
			log.Println(err)
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
		log.Fatalf("unknown type %v", string(t.MustBytes("")))
		return nil
	}
}

func build(ap core.ActionProvider, t typed.Typed, runner core.Runner) core.Check {
	name, ok := t.StringIf("name")
	if ok == false {
		log.Fatalf("missing name %v", string(t.MustBytes("")))
	}
	b := &Base{
		ap:        ap,
		name:      name,
		runner:    runner,
		recover:   t.Strings("recover"),
		snooze:    time.Second * time.Duration(t.IntOr("snooze", 0)),
		failCount: t.IntOr("failcount", 1),
	}
	if b.failCount <= 0 {
		b.failCount = 1
	}
	return b
}
