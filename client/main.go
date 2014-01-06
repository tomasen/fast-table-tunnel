
package main
import (
	"net"
	"log"
	"os"
	"syscall"
	"runtime/debug"
	ftunnel "bitbucket.org/Tomasen/fast-table-tunnel/src"
	gozd "bitbucket.org/PinIdea/zero-downtime-daemon"
)

const (
	_accept  = ":61080"
	_connect = "cdn.hk0.shooter.cn:61080"
)

const (
	BUFFER_MAXSIZE 	= 4096
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
			log.Panicln(err)
			return
		}
    
		ftunnel.Decrypt(buff[:n])
		
		_, err = client.Write(buff[:n])
		if err != nil {
			log.Panicln(err)
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
  
	server, err := net.Dial("tcp", _connect)
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
			log.Panicln(err)
			return
		}
    
		ftunnel.Encrypt(buff[:n])
		_, err = server.Write(buff[:n])
		if err != nil {
			log.Panicln(err)
			return
		}
	}
}

func handleListners(cl chan net.Listener) {
  
	for v := range cl {
		go func(l net.Listener){
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
	
	log.Println(os.TempDir())
	ctx  := gozd.Context{
		Hash:   "fast-tunnel",
		Command:"start",
		Maxfds: syscall.Rlimit{Cur:32677, Max:32677},
		User:   "www",
		Group:  "www",
		Logfile:"tunnel_daemon.log",
		Directives:map[string]gozd.Server{
			"client":gozd.Server{
				Network:"tcp",
				Address:_accept,
			},
		},
	}
  
	cl := make(chan net.Listener,1)
	go handleListners(cl)
	sig, err := gozd.Daemonize(ctx, cl) // returns channel that connects with daemon
	if err != nil {
		log.Println("error: ", err)
		return
	}
  
	// other initializations or config setting
	for s := range sig  {
		switch s {
		case syscall.SIGHUP, syscall.SIGUSR2:
			// do some custom jobs while reload/hotupdate
      
    
		case syscall.SIGTERM:
			// do some clean up and exit
			return
		}
	}
}

