package main

import (
	"flag"
	"github.com/karlseguin/thirdlaw"
)

var test = flag.Bool("test", false, "test the specified configuration file")
var config = flag.String("config", "config.json", "path to config file")

func main() {
	flag.Parse()
	if *test == true {
		thirdlaw.LoadConfig(*config)
		return
	}
	thirdlaw.Run(*config)
}
