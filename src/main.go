package ftunnel

import (
	"flag"
	"log"
)

var (
	optConf   = flag.String("conf", "", "config file or url")
	optHelp     = flag.Bool("h", false, "this help")
)

func usage() {
	log.Println("[command] -conf=[config file]")
	flag.PrintDefaults()
}

func Main() {
	
	// parse arguments
	flag.Parse()

	if *optHelp || len(*optConf) <= 0 {
		usage()
		return
	}

	// if config file is url
	// TODO: start config update daemon
	
	// if config is a file
	// TODO: just read and reload
	
}