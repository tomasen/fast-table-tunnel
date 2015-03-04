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
func NewTunnel() (tunnel *Tunnel, err error) {
	tunnel = &Tunnel{}
	//TODO use MTU set udp_pkg_size
	return
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
		buf := make([]byte, 0, this.tcp_buffer_size)
		err := tcp_conn.SetReadBuffer(cap(buf))
		if err != nil {
			logger.Warningln("loopReadTCP:SetReadBuffer:", err)
			break
		}
		_, err = tcp_conn.Read(buf)
		if err != nil {
			logger.Debugln("loopReadTCP:Read:", err)
			break
		}
		//TODO 压缩加密等预先处理tcp收到的数据包
		//count_retry, count_send := tcp_conn.Do(buf)
		//TODO 计算此次传输质量，计算tunnel质量，调节tcp_buffer_size
	}

	tcp_conn.Close()

	this.mu.Lock()
	defer this.mu.Unlock()

	delete(this.tcp_conns, tcp_conn.conn_id)
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
		//TODO go tcp_conn.Recv

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
