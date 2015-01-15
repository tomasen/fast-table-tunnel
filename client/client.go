package main

import (
	"ftunnel"
	"log"
	"net"
)

var (
	BUFFER_MAXSIZE = 4096
	tcp_addr       = "0.0.0.0:1357"
	udp_addr_str   = "71.19.157.50:2468"
	udp_addr       *net.UDPAddr
)

func main() {
	var err error
	udp_addr, err = net.ResolveUDPAddr("udp4", udp_addr_str)
	if err != nil {
		log.Fatal("main:ResolveUDPAddr:", err)
	}

	ln, err := net.Listen("tcp4", tcp_addr)
	if err != nil {
		log.Fatal("main:Listen:", err)
	}

	for {
		t_conn, err := ln.Accept()
		if err != nil {
			log.Fatal("main:Accept:", err)
		}

		laddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
		if err != nil {
			log.Fatal("main:ResolveUDPAddr:", err)
		}

		u_conn, err := net.ListenUDP("udp4", laddr)
		if err != nil {
			log.Fatal("main:ListenUDP:", err)
		}

		_, err = u_conn.WriteToUDP([]byte{}, udp_addr)
		if err != nil {
			log.Fatal("main:NewConnect:Write test:", err)
		}

		conn, err := ftunnel.NewConnect(t_conn, u_conn, udp_addr)
		if err != nil {
			log.Fatal("main:NewConnect:", err)
		}

		go conn.ListenAndServe()
	}
}
