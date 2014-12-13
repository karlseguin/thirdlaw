package main

import (
	"github.com/karlseguin/beats"
)

func main() {
	configPath := "config.json"
	beats.Start(configPath)
}
