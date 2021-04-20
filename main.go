package main

import (
	"flag"
	"log"
	"net"
)

const (
	BUFFER_MAXSIZE = 64 * 1024
)

func main() {
	addrListen := flag.String("s", "", "listen to ip:port")
	addrConn := flag.String("c", "", "connect to ip:port")
	flag.Parse()

	if len(*addrListen) == 0 || len(*addrConn) == 0 {
		flag.PrintDefaults()
		return
	}

	Serve(*addrListen, *addrConn)
}

func Serve(addrListen, addrConn string) {
	l, err := net.Listen("tcp", addrListen)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		connFromClient, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}

		go func() {
			connToServer, err := net.Dial("tcp", addrConn)
			if err != nil {
				log.Println("connect error:", err)
				return
			}

			go func() {
				defer connToServer.Close()
				defer connFromClient.Close()

				// from server to client
				buff := make([]byte, BUFFER_MAXSIZE)

				for {
					n, err := connToServer.Read(buff)
					if err != nil {
						log.Println("read from server error:", err)
						return
					}

					Flip(buff[:n])

					_, err = connFromClient.Write(buff[:n])
					if err != nil {
						log.Println("write to client error:", err)
						return
					}
				}
			}()

			go func() {
				defer connToServer.Close()
				defer connFromClient.Close()

				// from client to server
				buff := make([]byte, BUFFER_MAXSIZE)

				for {
					n, err := connFromClient.Read(buff)
					if err != nil {
						log.Println("read from client error", err)
						return
					}

					Flip(buff[:n])

					_, err = connToServer.Write(buff[:n])
					if err != nil {
						log.Println("write to server error:", err)
						return
					}
				}
			}()
		}()
	}
}

func Flip(buff []byte) {
	for k, v := range buff {
		buff[k] = ^v
	}
}
