package ftunnel

import (
	"net"
	"runtime/debug"
	"sync"

	"bitbucket.org/abotoo/gonohup"
)

const UDP_WRITER_MAXSIZE = 10000

var logger *gonohup.Logger

type Tunnel struct {
	r_udp_addr      string
	tcp_conns       map[uint16]*TCPConn
	tcp_buffer_size int32
	ch_udp_writer   chan []byte
	mu              sync.RWMutex
}

//TODO
func (this *Tunnel) Serve() {
	this.ch_udp_writer = make(chan []byte, UDP_WRITER_MAXSIZE)

	go this.loopReadUDP()
	go this.loopWriteUDP()

	logger.Info("tunnel was working, remote udp addr:", this.r_udp_addr)
}

//TODO
func (this *Tunnel) AppendTCPConn(tcp_conn *net.TCPConn) {
	this.mu.Lock()
	defer this.mu.Unlock()

	//TODO create TCPConn

	//go this.loopReadTCP(TCPConn)
}

//TODO
func (this *Tunnel) loopReadTCP(tcp_conn *TCPConn) {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("loopReadTCP recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
		}
	}()

	for {
		//TODO 读取数据，分割，存入buf，送入udp发送channel
		//TODO 创建定时，阻塞，超时或收到信号释放阻塞
		for {
			<-tcp_conn.ch_sig
			//TODO 要求重发,buf取出，送入udp发送channel
			//TODO Done，break loop
		}
	}

	//TODO 清理，从map删除TCPConn
}

//TODO
func (this *Tunnel) loopReadUDP() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("loopReadUDP recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
			go this.loopReadUDP()
		}
	}()

	for {
		//TODO 接收udp包，解析头部
		//TODO 查tcp conn id
		//TODO 信号包，送入tcp sig channel
		//TODO 组装完成的数据包，TCPConn.Write

		//TODO 超时未收到数据，发送keepAlive，超过阀值未收到回复，关闭本连接
	}
}

//TODO
func (this *Tunnel) loopWriteUDP() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("loopWriteUDP recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
			go this.loopWriteUDP()
		}
	}()

	for {
		<-this.ch_udp_writer
	}
}

//TODO
func (this *Tunnel) keepAlive() {

}
