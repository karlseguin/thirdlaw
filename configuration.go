package thirdlaw

import (
	"github.com/karlseguin/thirdlaw/actions"
	"github.com/karlseguin/thirdlaw/checks"
	"github.com/karlseguin/thirdlaw/core"
	"github.com/karlseguin/thirdlaw/outputs"
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"time"
)

type Configuration struct {
	frequency time.Duration
	checks    []core.Check
	onFailure []core.Output
	onSuccess []core.Output
	actions   map[string]core.Action
}

func loadConfig(path string) *Configuration {
	t, err := typed.JsonFile(path)
	if err != nil {
		panic(err)
	}

	chks := t.Objects("checks")
	if len(chks) == 0 {
		log.Println("WARN 0 checks configured")
	}

	onFailure, onSuccess := t.Object("outputs").Objects("failure"), t.Object("outputs").Objects("success")
	if len(onFailure) == 0 {
		log.Println("WARN 0 outputs configured for failure")
	}

	actns := t.Object("actions")

	c := &Configuration{
		checks:    make([]core.Check, len(chks)),
		onFailure: make([]core.Output, len(onFailure)),
		onSuccess: make([]core.Output, len(onSuccess)),
		frequency: time.Millisecond * time.Duration(t.IntOr("frequency", 10000)),
		actions:   make(map[string]core.Action, len(actns)),
	}

	for name, _ := range actns {
		c.actions[name] = actions.New(actns.Object(name))
	}
	for i, check := range chks {
		c.checks[i] = checks.New(c.actions, check)
	}
	for i, output := range onFailure {
		c.onFailure[i] = outputs.New(output)
	}
	for i, output := range onSuccess {
		c.onSuccess[i] = outputs.New(output)
	}

	return c
}
