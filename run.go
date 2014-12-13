package thirdlaw

import (
	"github.com/karlseguin/thirdlaw/core"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var reload bool

func Run(configPath string) {
	config := loadConfig(configPath)
	signals()
	for {
		time.Sleep(config.frequency)
		run(config)
		if reload {
			config = loadConfig(configPath)
			reload = false
		}
	}
}

func run(config *Configuration) {
	defer swallow()
	l := len(config.checks)

	success := make([]*core.Result, 0, l)
	failures := make([]*core.Result, 0, l)

	for i := 0; i < l; i++ {
		result := config.checks[i].Run()
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
	results := core.NewResults(list)
	for _, output := range outputs {
		output.Process(results)
	}
}

func swallow() {
	if err := recover(); err != nil {
		log.Println("unhandled error", err)
		time.Sleep(time.Second * 5)
	}
}

func signals() {
	sigusr2 := make(chan os.Signal, 1)
	signal.Notify(sigusr2, syscall.SIGUSR2)
	go func() {
		for {
			<-sigusr2
			reload = true
		}
	}()
}
