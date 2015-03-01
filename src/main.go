// entry point
package ftunnel

import (
	"flag"
	"log"
)

var (
	optConf = flag.String("conf", "", "config file or url")
	optHelp = flag.Bool("h", false, "this help")
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

	// TODO: fork

	var sp Supervisor
	err := sp.Load(*optConf)
	if err != nil {
		log.Println("E(main.Supervisor.Source):", err)
		return
	}

	c := make(chan bool, 1)
	for _ = range c {}
}
