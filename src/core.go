// services control center
package ftunnel

import (
	"log"
	"net"
	"math/rand"
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

func (co *Core) PushPackedDataToNextNode(b []byte) {
	// TODO: sort by score
	// pick one of top 3 randomly
	// PushNextPackedData
	l := len(co.Nodes)
	n := l - 1
	if n > 3 {
		n = 3
	}
	last := -1
	for i := n; i >= 0; i-- {
		r := rand.Intn(l)
		if co.Nodes[r].Identity == _nodeId {
			i++
			continue
		}
		if last >= 0 {
			if co.Nodes[r].score < co.Nodes[last].score {
				last = r
			}
		} else {
			last = r
		}
	}

	for _, n := range co.Nodes {
		if n.Identity != _nodeId {
			n.PushPacket(b)
		}
	}
}
