package ftunnel

import (
	"log"
	"net"
)

type Service struct {
	Name          string
	Network       string // "tcp", "udp", "ipv4"
	Address       string // ":9000"
	OutboundGroup int
	InboundGroup  int
	DstIp         string
	DstPort       string
	tcp_l         net.Listener
	Nodes         *[]Node
}

func (s *Service) Start() {
	// TODO: udp gre raw_ipv4
	var err error
	s.tcp_l, err = net.Listen("tcp", s.Address)
	if err != nil {
		log.Println("E(service.StartListen):", err)
		return
	}

	for {
		conn, err := s.tcp_l.Accept()
		if err != nil {
			// handle error
			log.Println("N(service.Accept):", err)
			continue
		}

		go func(c net.Conn) {
			var b []byte

			// TODO: build connection by
			// direct connection (if myself is in the outbounf group)
			// and send conn Packet to next Hop
			// connid := ConnId()

			for {
				_, err := c.Read(b)
				if err != nil {
					log.Println("E(service.Serv):", err)
					break
				}
				// TODO: send to other nodes, smartly?
				// if this node is inbound or outbound
				// if this is outbound send one to dst
			}
		}(conn)
	}
}

func (s *Service) Stop() {
	// TODO: udp gre raw_ipv4
	s.tcp_l.Close()
}
