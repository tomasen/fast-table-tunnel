// entry point
package ftunnel

import (
	"flag"
	"log"
)

func usage() {
	log.Println("[command] -conf=[config file]")
	flag.PrintDefaults()
}

func MainSimple() {

	optServerMode := flag.Bool("s", false, "server mode")
	optConnectTo := flag.String("conn", "", "destination protocol:ip:port")
	optListenTo  := flag.String("listen", "", "listen to protocol:ip:port")
	optHelp := flag.Bool("h", false, "this help")

	// parse arguments
	flag.Parse()

	if *optHelp || len(*optConnectTo) <= 0 || len(*optListenTo) <= 0 {
		usage()
		return
	}

	// TODO: fork
	if *optServerMode != false	{
		
	}
	exit := make(chan bool, 1)
	<- exit
}


func MainV2() {

	optConf := flag.String("conf", "", "config file or url")
	optHelp := flag.Bool("h", false, "this help")

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

	exit := make(chan bool, 1)
	<- exit
}
