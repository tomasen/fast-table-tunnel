// services control center
package ftunnel

type Node struct {
	Groups   []int
	Ip       string
	IsMyself bool
}

type Core struct {
	Nodes          []Node
	Services       []Service
	BinaryUrl      string
	BinaryCheckSum []byte
}

func (co *Core) Start() {
	// TODO: determine which node is this node
	// query with other nodes or http://ifconfig.me/ip ipinfo.io/ip
	myip, err := myip()
	if err != nil {
		for _, s := range co.Nodes {
			// TODO: query other nodes for my ip
		}
	}
	for _, s := range co.Nodes {
		s.IsMyself = (s.Ip == myip)
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
