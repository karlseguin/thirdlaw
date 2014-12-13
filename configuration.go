package beats

import (
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"time"
)

type Configuration struct {
	frequency time.Duration
	checks    []Check
}

func loadConfig(path string) *Configuration {
	t, err := typed.JsonFile(path)
	if err != nil {
		panic(err)
	}

	checks := t.Objects("checks")
	if len(checks) == 0 {
		log.Println("0 checks configured")
	}
	c := &Configuration{
		checks:    make([]Check, len(checks)),
		frequency: time.Millisecond * time.Duration(t.IntOr("frequency", 10000)),
	}
	for i, check := range checks {
		c.checks[i] = CheckFactory(check)
	}
	return c
}
