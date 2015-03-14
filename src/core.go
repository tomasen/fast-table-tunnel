// services control center
package ftunnel

import (
	"log"
	"net"
)

var (
//	_core Core
)

type Core struct {
	Nodes          []Node
	Services       []Service
	BinaryUrl      string
	BinaryCheckSum []byte
	listener       net.Listener
}

func (co *Core) Start() {
	// determine which node is this node
	myip := ip()
	for _, nd := range co.Nodes {
		if myip == nd.Ip {
			nd.Identity = _nodeId
			co.StartListen(nd.Port)
		} else {
			// check identity of other node
			go nd.Connect()
		}
	}

	// start all services
	for _, s := range co.Services {
		// TODO: Elect Next Hop
		s.co = co
		s.Start()
	}
}

func (co *Core) Stop() {
	// close all nodes connections
	for _, n := range co.Nodes {
		n.Close()
	}

	co.listener.Close()
	// close all services
	for _, s := range co.Services {
		//if is tcplistener
		s.Stop()
	}
}

func (co *Core) StartListen(port string) {
	var err error
	co.listener, err = net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("E(core.StartListen):", err)
		return
	}

	for {
		conn, err := co.listener.Accept()
		if err != nil {
			// handle error
			log.Println("N(core.Accept):", err)
			continue
		}
		tr := NewTransporter(conn)
		go tr.ServConnection()
	}
}

func (co *Core) NodeDoesBelongGroup(group int, nodeid uint64) bool {
	for _, n := range co.Nodes {
		if n.Identity == nodeid {
			for _, g := range n.Groups {
				if g == group {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (co *Core) PushPacketToAllNodes(b []byte) {
	// TODO: 
	for _, n := range co.Nodes {
		if n.Identity != _nodeId {
			n.PushPacket(b)
		}
	}
}