package main

import (
	gozd "bitbucket.org/PinIdea/zero-downtime-daemon"
	ftunnel "bitbucket.org/Tomasen/fast-table-tunnel/src"
	"flag"
	"log"
	"net"
	"os"
	"syscall"
)

var (
	optAccept   = flag.String("s", "", "listen to ip:port")
	optConnect  = flag.String("c", "", "connect to ip:port")
	optServerID = flag.String("id", "fasttunnel", "connect to ip:port")
	optLogfile  = flag.String("log", "", "log filepath")
	optHelp     = flag.Bool("h", false, "this help")
)

func usage() {
	log.Println("[command] -conf=[config file]")
	flag.PrintDefaults()
}

func main() {

	// parse arguments
	flag.Parse()

	if *optHelp || len(*optServerID) <= 0 || len(*optAccept) <= 0 || len(*optConnect) <= 0 {
		usage()
		return
	}

	log.Println(os.TempDir())
	ctx := gozd.Context{
		Hash:    *optServerID,
		Command: "start",
		Maxfds:  syscall.Rlimit{Cur: 32677, Max: 32677},
		User:    "nobody",
		Group:   "nobody",
		Logfile: "ftunnel.log",
		Directives: map[string]gozd.Server{
			"client": gozd.Server{
				Network: "tcp",
				Address: *optAccept,
			},
		},
	}

	cl := make(chan net.Listener, 1)
	go ftunnel.HandleListners(cl, *optConnect)
	sig, err := gozd.Daemonize(ctx, cl) // returns channel that connects with daemon
	if err != nil {
		log.Println("error: ", err)
		return
	}

	// other initializations or config setting
	for s := range sig {
		switch s {
		case syscall.SIGHUP, syscall.SIGUSR2:
			// do some custom jobs while reload/hotupdate

		case syscall.SIGTERM:
			// do some clean up and exit
			return
		}
	}
}
