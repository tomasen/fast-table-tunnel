package ftunnel

type Service struct {
	Name          string
	Network       string // "tcp", "udp", "ipv4"
	Address       string // ":9000"
	OutboundGroup int
	InboundGroup  int
	DstIp         string
	DstPort       string
}

func (s *Service) Start() {
	// TODO:
}

func (s *Service) Stop() {
	// TODO:
}
