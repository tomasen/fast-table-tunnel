package ftunnel

import "net"

var (
	udp_conn_in = make(map[string]*net.UDPConn)
)

//TODO
func ListenAndServeUDP(l_udp_addr, r_tcp_addr string) (err error) {
	//TODO dial r_tcp_addr测试连接
	//TODO listen l_udp_addr
	//TODO 检查重复，r_udp_addr已经存在，则丢弃新请求
	//TODO dial r_udp_addr, 新建tunnel，设置udp_pkg_size
	//go tunnel.Serve()
	return
}

//TODO
func dialTCP() {

}
