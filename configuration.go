package thirdlaw

import (
	"github.com/karlseguin/thirdlaw/actions"
	"github.com/karlseguin/thirdlaw/checks"
	"github.com/karlseguin/thirdlaw/core"
	"github.com/karlseguin/thirdlaw/outputs"
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"io/ioutil"
	"time"
)

type Configuration struct {
	frequency time.Duration
	checks    []core.Check
	onFailure []core.Output
	onSuccess []core.Output
	actions   map[string]core.Action
}

func (c *Configuration) GetAction(name string) core.Action {
	return c.actions[name]
}

func loadConfig(path string) *Configuration {
	t, err := typed.JsonFile(path)
	if err != nil {
		panic(err)
	}
	onFailure, onSuccess := t.Object("outputs").Objects("failure"), t.Object("outputs").Objects("success")
	if len(onFailure) == 0 {
		log.Println("WARN 0 outputs configured for failure")
	}

	c := &Configuration{
		checks:    make([]core.Check, 0, 20),
		actions:   make(map[string]core.Action),
		onFailure: make([]core.Output, len(onFailure)),
		onSuccess: make([]core.Output, len(onSuccess)),
		frequency: time.Millisecond * time.Duration(t.IntOr("frequency", 10000)),
	}
	for i, output := range onFailure {
		c.onFailure[i] = outputs.New(output)
	}
	for i, output := range onSuccess {
		c.onSuccess[i] = outputs.New(output)
	}
	loadOne(c, t)
	if include, ok := t.StringIf("include"); ok {
		files, err := ioutil.ReadDir(include)
		if err != nil {
			panic(err)
		}
		if include[len(include)-1] != '/' {
			include += "/"
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			data, err := ioutil.ReadFile(include + file.Name())
			if err != nil {
				panic(err)
			}
			t, err = typed.Json(data)
			if err != nil {
				panic(err)
			}
			loadOne(c, t)
		}
	}
	return c
}

func loadOne(c *Configuration, t typed.Typed) {
	for _, check := range t.Objects("checks") {
		c.checks = append(c.checks, checks.New(c, check))
	}
	actns := t.Object("actions")
	for name, _ := range actns {
		c.actions[name] = actions.New(name, actns.Object(name))
	}

	if check, ok := t.ObjectIf("check"); ok {
		c.checks = append(c.checks, checks.New(c, check))
	}
}
