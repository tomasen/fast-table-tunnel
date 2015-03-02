package ftunnel

import "net"

var (
	udp_conn_in = make(map[string]*net.UDPConn)
)

//TODO
func ListenAndServeUDP(l_udp_addr, r_tcp_addr string) (err error) {
	//TODO:dial r_tcp_addr获得连接
	//TODO:listen l_udp_addr
	//TODO:检查重复，udp_conn_in已经存在，则丢弃新连接
	return
}

//TODO
func dialTCP() {

}
