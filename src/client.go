package ftunnel

var (
	MTU = 1492
)

//TODO
func ListenAndServeTCP(l_tcp_addr, r_udp_addr string) (err error) {
	getMTU(r_udp_addr)
	dialUDP(r_udp_addr)
	//TODO listen l_tcp_addr
	return
}

//TODO
func dialUDP(r_udp_addr string) {
	//TODO dial r_udp_addr获得连接，通知udp_pkg_size
	//TODO 监听udp连接，获取对方的回复的此次将要tunnel用的新r_udp_addr
	//TODO 新建tunnel，设置udp_pkg_size
	//go tunnel.Serve()
}

//TODO 探测MTU
func getMTU(r_udp_addr string) {
	//TODO 根据MTU算出udp_pkg_size
}
