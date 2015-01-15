package main

import (
	"ftunnel"
	"log"
	"net"
)

var (
	BUFFER_MAXSIZE = 4096
	tcp_addr       = "127.0.0.1:55552"
	udp_addr_str   = "0.0.0.0:2468"
	udp_addr       *net.UDPAddr
)

func main() {
	var err error
	udp_addr, err = net.ResolveUDPAddr("udp4", udp_addr_str)
	if err != nil {
		log.Fatal("main:ResolveUDPAddr:", err)
	}

	udp_listener, err := net.ListenUDP("udp4", udp_addr)
	if err != nil {
		log.Fatal("main:Listen:", err)
	}
	//TODO:need lock map to manager
	connMap := make(map[string]*ftunnel.Connect)
	for {
		buff := make([]byte, BUFFER_MAXSIZE)
		n, u_addr, err := udp_listener.ReadFromUDP(buff)
		if err != nil {
			log.Fatal("main:ReadFromUDP:", err)
		}
		if n < 5 {
			continue
		}
		u_addr_str := u_addr.String()
		if c, exist := connMap[u_addr_str]; exist {
			go c.Recv(buff[:n])
		} else if buff[4] != ftunnel.SIGNAL_RESEND {
			t_conn, err := net.Dial("tcp4", tcp_addr)
			if err != nil {
				log.Fatal("main:Dial tcp:", err)
			}

			conn, err := ftunnel.NewConnect(t_conn, udp_listener, u_addr)
			if err != nil {
				log.Fatal("main:NewConnect:", err)
			}
			connMap[u_addr_str] = conn
			go conn.Serve()
			go conn.Recv(buff[:n])
		}
	}
}
