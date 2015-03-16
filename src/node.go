// Node operations
package ftunnel

import (
	"log"
	"math"
	"net"
	"time"
	"sort"
)

type Node struct {
	Groups   []int
	Ip       string
	Port     string
	Identity uint64
	score    int64 // lantency
	tr       *Transporter
	co       *Core
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

		// TODO: handle node removal?
	}

	// keep ping to rate the score of node
	for {
		nd.score = nd.tr.Ping()
		if math.MaxInt64 == nd.score {
			// Reconnect
			go nd.Connect()
			break
		}
		sort.Sort(nodes(nd.co.Nodes))
		time.Sleep(1 * time.Second)
		// Handle node removal
		if nd.tr == nil {
			return
		}
	}
}

func (nd *Node) Close() {
	nd.tr.Close()
	nd.tr = nil
}

func (nd *Node) PushPacket(b []byte) {
	if nd.tr == nil {
		nd.tr.WritePacketBytes(b)
	}
}

type nodes []Node

func (a nodes) Len() int           { return len(a) }
func (a nodes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nodes) Less(i, j int) bool { return a[i].score < a[j].score }
