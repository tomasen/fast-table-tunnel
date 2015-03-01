// services control center
package ftunnel

type Node struct {
	Groups   []int
	Ip       string
	Identity uint64
}

type Core struct {
	Nodes          []Node
	Services       []Service
	BinaryUrl      string
	BinaryCheckSum []byte
}

func (co *Core) Start() {
	// determine which node is this node
	myip := ip()
	for _, s := range co.Nodes {
		if myip == s.Ip {
			s.Identity = _nodeId
		} else {
			// TODO: check identity of other node

		}
	}

	// start all services
	for _, s := range co.Services {
		s.Start()
	}
}

func (co *Core) Stop() {
	// close all services
	for _, s := range co.Services {
		//if is tcplistener
		s.Stop()
	}
}
