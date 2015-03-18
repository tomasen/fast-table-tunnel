package ftunnel

import (
	"log"
	"net"
	"runtime/debug"
)

const (
	BUFFER_MAXSIZE = 64 * 1024
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

		Decrypt(buff[:n])

		_, err = client.Write(buff[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func serveTCP(client net.Conn, adr string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in", r, ":")
			log.Println(string(debug.Stack()))
		}
	}()

	defer client.Close()

	server, err := net.Dial("tcp", adr)
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

		Encrypt(buff[:n])
		_, err = server.Write(buff[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func HandleListners(cl chan net.Listener, adr string) {

	for v := range cl {
		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				if err != nil {
					// gozd.ErrorAlreadyStopped may occur when shutdown/reload
					log.Println("accept error: ", err)
					break
				}

				go serveTCP(conn, adr)
			}
		}(v)
	}
}
