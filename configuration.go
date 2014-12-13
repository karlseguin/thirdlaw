package beats

import (
	"gopkg.in/karlseguin/typed.v1"
	"github.com/karlseguin/beats/core"
	"github.com/karlseguin/beats/checks"
	"log"
	"time"
)

type Configuration struct {
	frequency time.Duration
	checks    []core.Check
}

func loadConfig(path string) *Configuration {
	t, err := typed.JsonFile(path)
	if err != nil {
		panic(err)
	}

	chks := t.Objects("checks")
	if len(chks) == 0 {
		log.Println("0 checks configured")
	}
	c := &Configuration{
		checks:    make([]core.Check, len(chks)),
		frequency: time.Millisecond * time.Duration(t.IntOr("frequency", 10000)),
	}
	for i, check := range chks {
		c.checks[i] = checks.New(check)
	}
	return c
}
