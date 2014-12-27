package thirdlaw

import (
	"github.com/karlseguin/thirdlaw/core"
	"log"
	"time"
)

func Run(configPath string) {
	config := LoadConfig(configPath)
	for {
		time.Sleep(config.frequency)
		run(config)
	}
}

func run(config *Configuration) {
	defer swallow()
	l := len(config.checks)

	success := make([]*core.Result, 0, l)
	failures := make([]*core.Result, 0, l)

	for i := 0; i < l; i++ {
		result := config.checks[i].Run()
		if result == nil {
			continue
		}
		if result.Ok {
			success = append(success, result)
		} else {
			failures = append(failures, result)
		}
	}

	outputs, list := config.onSuccess, success
	if len(failures) > 0 {
		outputs, list = config.onFailure, failures
	}
	if len(list) > 0 {
		results := core.NewResults(list)
		for _, output := range outputs {
			output.Process(results)
		}
	}
}

func swallow() {
	if err := recover(); err != nil {
		log.Println("unhandled error", err)
		time.Sleep(time.Second * 5)
	}
}
