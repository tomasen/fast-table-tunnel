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
	DstAddress    string // IP:Port
	tcp_l         net.Listener
	co            *Core
}

func (s *Service) Start() {
	if !s.co.NodeDoesBelongGroup(s.InboundGroup, _nodeId) {
		return
	}
	// TODO: listen to udp gre raw_ipv4
	var err error
	s.tcp_l, err = net.Listen("tcp", s.Address)
	if err != nil {
		log.Println("E(service.StartListen):", err)
		return
	}

	outbound := s.co.NodeDoesBelongGroup(s.OutboundGroup, _nodeId)

	for {
		conn, err := s.tcp_l.Accept()
		if err != nil {
			// handle error
			log.Println("N(service.Accept):", err)
			continue
		}

		go func(c net.Conn, d bool, connid uint64) {
			var b []byte

			// build connection by
			if d {
				// TODO: direct connection (if myself is in the outbound group)
				// which might shouldn't be happenning
			} else {
				// TODO: send conn Packet to next Hop
				b := BuildConnPacket("tcp", s.DstAddress)
				s.co.PushPacketToAllNodes(b)
			}

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
		}(conn, outbound, ConnId())
	}
}

func (s *Service) Stop() {
	// TODO: udp gre raw_ipv4
	if s.tcp_l != nil {
		s.tcp_l.Close()
	}
}
