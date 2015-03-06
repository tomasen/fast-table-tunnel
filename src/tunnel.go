package ftunnel

import (
	"net"
	"runtime/debug"
	"sync"
	"time"

	"bitbucket.org/abotoo/gonohup"
)

const (
	UDP_WRITER_MAXSIZE = 10000
	UDP_READ_TIMEOUT   = 60 //Second
	ALIVE_MAX          = uint8(5)
)

var (
	logger                  *gonohup.Logger
	udp_pkg_header_size     int
	default_tcp_buffer_size int32
)

type Tunnel struct {
	r_udp_addr      net.Addr
	l_udp_addr      net.Addr
	tcp_conns       map[uint16]*TCPConn
	tcp_buffer_size int16
	udp_pkg_size    int16
	ch_udp_writer   chan []byte
	mu              sync.RWMutex
	alive_counter   uint8
}

//TODO
func NewTunnel(l_udp_addr, r_udp_addr net.Addr) (tunnel *Tunnel, err error) {
	tunnel = &Tunnel{l_udp_addr: l_udp_addr, r_udp_addr: r_udp_addr, alive_counter: 0}
	tunnel.tcp_conns = make(map[uint16]*TCPConn)
	tunnel.tcp_buffer_size = default_tcp_buffer_size
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
		if this.tcp_buffer_size < this.udp_pkg_size {
			logger.Warningln("loopReadTCP:tcp_buffer_size too small:", this.tcp_buffer_size)
			break
		}
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

	var laddr *net.UDPAddr
	var err error
	if this.l_udp_addr != nil {
		laddr, err = net.ResolveUDPAddr(this.l_udp_addr.Network(), this.l_udp_addr.String())
		if err != nil {
			logger.Warningln("tunnel:ResolveUDPAddr:", err)
			logger.Debugln(this.l_udp_addr)
			return
		}
	} else {
		logger.Warningln("tunnel:l_udp_addr is nil")
		return
	}
	udp_conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		logger.Warningln("tunnel:ListenUDP:", err)
		return
	}

	for {
		buf := make([]byte, 0, this.udp_pkg_size)
		udp_conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(UDP_READ_TIMEOUT)))
		n, err := udp_conn.Read(buf)
		if err != nil {
			if n_err, flag := err.(net.Error); flag {
				if n_err.Timeout() {
					this.alive_counter += 1
					if this.alive_counter > ALIVE_MAX {
						logger.Warningln("tunnel:keepAlive over max")
						break
					}
					this.keepAlive()
					continue
				}
			}
			logger.Warningln("tunnel:udp_conn read:", err)
			return
		}
		logger.Debugln("tunnel:udp_conn read:", n)

		if len(buf) < udp_pkg_header_size {
			logger.Debugf("tunnel:read too small pkg: %v\n", buf)
			continue
		}

		header := NewExtraHeaderByPkg(buf[:udp_pkg_header_size])
		if header == nil {
			logger.Debugf("tunnel:get header failed: %v\n", buf[:udp_pkg_header_size])
			continue
		}

		//TODO 判断是否tunnel的信号包，keepAlive重置或要求设定udp_pkg_size
		//this.alive_counter=0

		var tcp_conn *TCPConn
		var exist bool
		conn_id := header.GetTCPConnId()

		this.mu.RLock()
		tcp_conn, exist = this.tcp_conns[conn_id]
		this.mu.RUnlock()

		if exist {
			go tcp_conn.Recv(buf)
		} else {
			logger.Warningln("tunnel:can't found tcp conn by:", conn_id)
		}
	}
}

//loopWriteUDP
func (this *Tunnel) loopWriteUDP() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("loopWriteUDP recovery: %v\n", r)
			logger.Debugf("%v\n", debug.Stack())
			go this.loopWriteUDP()
		}
	}()

	var laddr, raddr *net.UDPAddr
	var err error
	if this.l_udp_addr != nil {
		laddr, err = net.ResolveUDPAddr(this.l_udp_addr.Network(), this.l_udp_addr.String())
		if err != nil {
			logger.Warningln("tunnel:ResolveUDPAddr:", err)
			logger.Debugln(this.l_udp_addr)
			return
		}
	}
	if this.r_udp_addr != nil {
		raddr, err = net.ResolveUDPAddr(this.r_udp_addr.Network(), this.r_udp_addr.String())
		if err != nil {
			logger.Warningln("tunnel:ResolveUDPAddr:", err)
			logger.Debugln(this.r_udp_addr)
			return
		}
	} else {
		logger.Warningln("tunnel:r_udp_addr is nil")
		return
	}
	udp_conn, err := net.DialUDP("udp4", laddr, raddr)
	if err != nil {
		logger.Warningln("tunnel:DialUDP:", err)
		return
	}

	for {
		n, err := udp_conn.Write(<-this.ch_udp_writer)
		if err != nil {
			logger.Warningln("tunnel:udp_conn write:", err)
			return
		}
		logger.Debugln("tunnel:udp_conn write:", n)
	}
}

//TODO
func (this *Tunnel) keepAlive() {

}
