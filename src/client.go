package ftunnel

var (
	MTU          = 1492
	udp_pkg_size int
)

//TODO
func ListenAndServeTCP(l_tcp_addr, r_udp_addr string) (err error) {
	//TODO dial r_udp_addr获得连接
	//TODO listen l_tcp_addr
	return
}

//TODO
func dialUDP() {

}

//TODO 探测MTU
func getMTU() {
	//TODO 根据MTU算出最小udp_pkg_size
}
