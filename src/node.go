// Node operations
package ftunnel

import (
	"log"
	"net"
	"time"
)

type Node struct {
	Groups   []int
	Ip       string
	Port     string
	Identity uint64
	score    int64 // lantency
	tr       *Transporter
}

func (nd *Node) Connect() {
	for {
		conn, err := net.Dial("tcp", nd.Ip+":"+nd.Port)
		if err == nil {
			// send query packet to this node that ask for identity
			nd.tr = NewTransporter(conn)
			nd.Identity = nd.tr.QueryIdentity()
			break
		}

		log.Println("E(node.CheckIdentity):", err)
		time.Sleep(3 * time.Second)
		// TODO: proper handle node removal
	}

	// keep ping to rate the score of node
	for {
		nd.score = nd.tr.Ping()
		time.Sleep(1 * time.Second)
	}
}
