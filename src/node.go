// Node operations
package ftunnel

import (
	"log"
	"net"
)

type Node struct {
	Groups   []int
	Ip       string
	Port     string
	Identity uint64
	conn     net.Conn
}

func (nd *Node) CheckIdentity() (err error) {
	for {
		nd.conn, err = net.Dial("tcp", nd.Ip+":"+nd.Port)
		if err != nil {
			log.Println("E(node.CheckIdentity):", err)
			continue
		}
		// TODO: send query packet to this node that ask for identity

	}
}

// TODO: keep ping to rate the score of node
