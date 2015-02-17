package main

import (
	gozd "bitbucket.org/PinIdea/zero-downtime-daemon"
	ftunnel "bitbucket.org/Tomasen/fast-table-tunnel/src"
	"flag"
	"log"
	"net"
	"os"
	"runtime/debug"
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

const (
	BUFFER_MAXSIZE = 4096
)

func handleServer(client, server net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in", r, ":")
			log.Println(string(debug.Stack()))
		}
	}()

	defer client.Close()
	defer server.Close()

	buff := make([]byte, BUFFER_MAXSIZE)

	for {
		n, err := server.Read(buff)
		if err != nil {
			log.Println(err)
			return
		}

		ftunnel.Decrypt(buff[:n])

		_, err = client.Write(buff[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func serveTCP(client net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in", r, ":")
			log.Println(string(debug.Stack()))
		}
	}()

	defer client.Close()

	server, err := net.Dial("tcp", *optConnect)
	if err != nil {
		// handle error
		log.Panicln(err)
	}

	go handleServer(client, server)

	defer server.Close()

	buff := make([]byte, BUFFER_MAXSIZE)

	for {
		n, err := client.Read(buff)
		if err != nil {
			log.Println(err)
			return
		}

		ftunnel.Encrypt(buff[:n])
		_, err = server.Write(buff[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func handleListners(cl chan net.Listener) {

	for v := range cl {
		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				if err != nil {
					// gozd.ErrorAlreadyStopped may occur when shutdown/reload
					log.Println("accept error: ", err)
					break
				}

				go serveTCP(conn)
			}
		}(v)
	}
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
		User:    "www",
		Group:   "www",
		Logfile: "tunnel_daemon.log",
		Directives: map[string]gozd.Server{
			"client": gozd.Server{
				Network: "tcp",
				Address: *optAccept,
			},
		},
	}

	cl := make(chan net.Listener, 1)
	go handleListners(cl)
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
