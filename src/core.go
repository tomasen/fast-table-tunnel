// services control center
package ftunnel

type Core struct {
	Nodes          []Node
	Services       []Service
	BinaryUrl      string
	BinaryCheckSum []byte
}

func (co *Core) Start() {
	// determine which node is this node
	myip := ip()
	for _, nd := range co.Nodes {
		if myip == nd.Ip {
			nd.Identity = _nodeId
		} else {
			// check identity of other node
			go nd.CheckIdentity()
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
