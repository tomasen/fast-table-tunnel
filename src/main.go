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
	c := make(chan []byte, 1)
	var cfg Config
	err := cfg.Load(*optConf, c)
	if err != nil {
		log.Println("ERR:", err)
		return
	}
}
