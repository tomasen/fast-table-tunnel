package ftunnel

import (
	"net"
	"runtime/debug"
)

type TCPConn struct {
	tcp_conn  *net.TCPConn
	conn_id   uint16
	ch_sig    chan uint8
	ch_udp_in []byte
	buf       []byte
}

//TODO
func (this *TCPConn) Write(data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("TCPConn Write recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
		}
	}()

	return
}

//TODO
func (this *TCPConn) Close() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("TCPConn Close recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
		}
	}()

	close(this.ch_sig)
}

//TODO
func (this *TCPConn) zuzhuang() {
	//TODO 超时要求重发，收到要求丢弃的信号，丢弃本次拼接
}
