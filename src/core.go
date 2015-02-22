package ftunnel


var (
	_core Core
)

type Node struct {
	Groups []int
	Ip     string
}

type Core struct {
	Nodes          []Node
	Services       []Service
	BinaryUrl      string
	BinaryCheckSum []byte
}

func (co *Core) Start() {
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


