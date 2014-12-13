package beats

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/karlseguin/beats/core"
	"time"
)

var reload bool

func Start(configPath string) {
	config := loadConfig(configPath)
	signals()
	for {
		time.Sleep(config.frequency)
		beat(config)
		if reload {
			config = loadConfig(configPath)
			reload = false
		}
	}
}

func beat(config *Configuration) {
	defer swallow()
	l := len(config.checks)
	results := make([]*core.Result, l)
	for i := 0; i < l; i++ {
		results[i] = config.checks[i].Run()
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
