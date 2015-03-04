package ftunnel

import (
	"errors"
	"net"
	"runtime/debug"
	"time"
)

var (
	TCP_READ_TIMEOUT = 180 //Second
)

type TCPConn struct {
	conn            *net.TCPConn
	conn_id         uint16
	ch_sig_pkg      chan []byte
	ch_udp_in       chan []byte
	ch_udp_out      chan []byte //tunnel's channel, don't be close
	udp_out_buf     map[uint8][]byte
	udp_in_buf      map[uint8][]byte
	tcp_buffer_size int
}

//TODO
func NewTCPConn(conn_id uint16, conn *net.TCPConn, ch_udp_out chan []byte) (tcp_conn *TCPConn, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("NewTCPConn recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
			err = errors.New("NewTCPConn panic")
		}
	}()
	tcp_conn = &TCPConn{conn_id: conn_id, conn: conn, ch_udp_out: ch_udp_out}
	tcp_conn.ch_sig_pkg = make(chan []byte)
	tcp_conn.ch_udp_in = make(chan []byte)
	tcp_conn.udp_in_buf = make(map[uint8][]byte)
	tcp_conn.udp_out_buf = make(map[uint8][]byte)

	go tcp_conn.join()
	return
}

//Read
func (this *TCPConn) Read(buf []byte) (int, error) {
	this.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(TCP_READ_TIMEOUT)))
	return this.conn.Read(buf)
}

//SetReadBuffer
func (this *TCPConn) SetReadBuffer(tcp_buffer_size int) error {
	if tcp_buffer_size < udp_pkg_size {
		return errors.New("SetReadBuffer too small")
	} else if this.tcp_buffer_size != tcp_buffer_size {
		this.tcp_buffer_size = tcp_buffer_size
		return this.conn.SetReadBuffer(tcp_buffer_size)
	}
	return nil
}

//TODO
func (this *TCPConn) Do(data []byte) (count_retry, count_send int) {
	//TODO 分割，存入buf，送入udp发送channel
	//TODO 创建定时，阻塞，超时或收到信号释放阻塞
	for {
		<-this.ch_sig_pkg
		//TODO 超时重发,buf取出游标后的包，送入udp发送channel
		//TODO 收到确认游标包，移动游标
		//TODO 收到DONE，break loop
	}

}

//TODO
func (this *TCPConn) Recv(data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("TCPConn Write recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
		}
	}()

	//TODO check pkg header, size

	//TODO 信号包
	this.ch_sig_pkg <- data

	//TODO 数据包
	this.ch_udp_in <- data
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

	//TODO check this.udp_out_buf, then close ch_sig_pkg
	close(this.ch_sig_pkg)

	//TODO check this.udp_in_buf, then close ch_udp_in
	close(this.ch_udp_in)
}

//TODO goroutine, join packet
func (this *TCPConn) join() {
	//TODO 组装
	<-this.ch_udp_in
	//TODO 返回游标确认
	//this.ch_udp_out<-
}
