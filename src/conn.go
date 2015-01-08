package ftunnel

import (
	"log"
	"net"
	"sync"
	"time"
)

var (
	BUFFER_MAXSIZE = 4096
	CONN_TIMEOUT   = 15 //seconds tcp and udp connect timeout
)

//TCP and UDP connect
type Connect struct {
	in       map[uint16]*UDPPackage
	out      map[uint16]*UDPPackage
	tcp_conn net.Conn
	udp_conn *net.UDPConn
	cur      uint16
	locker   sync.Mutex
}

//NewConnect
func NewConnect(tcp_conn net.Conn, udp_conn *net.UDPConn) (*Connect, error) {
	var err error
	c := &Connect{tcp_conn: tcp_conn, udp_conn: udp_conn, cur: 0}
	c.in = make(map[uint16]*UDPPackage)
	c.out = make(map[uint16]*UDPPackage)
	//TODO:do some check
	return c, err
}

//create new UDPPackage
func (this *Connect) newUDPPackage(isOut bool, conn_id uint16) *UDPPackage {
	this.locker.Lock()
	defer this.locker.Unlock()

	if isOut {
		for {
			this.cur = this.cur + 1
			if _, exist := this.out[this.cur]; !exist {
				udpp := NewUDPPackage(this.tcp_conn, this.udp_conn, this.cur)
				this.out[this.cur] = udpp
				return udpp
			}
		}
	} else {
		udpp := NewUDPPackage(this.tcp_conn, this.udp_conn, conn_id)
		this.in[conn_id] = udpp
		return udpp
	}
	return nil
}

//Serve
func (this *Connect) Serve() {
	defer this.tcp_conn.Close()
	buff := make([]byte, BUFFER_MAXSIZE)
	for {
		this.tcp_conn.SetDeadline(time.Now().Add(time.Second * time.Duration(CONN_TIMEOUT)))
		n, err := this.tcp_conn.Read(buff)
		if err != nil {
			log.Println("Connect:Serve:Read:", err)
			return
		}

		err = this.newUDPPackage(true, 0).Send(buff[:n])
		if err != nil {
			log.Println("Connect:Serve:Send:", err)
			return
		}
	}
	return
}

//ListenUDP
func (this *Connect) ListenUDP(udp_conn *net.UDPConn) {
	defer udp_conn.Close()
	buff := make([]byte, BUFFER_MAXSIZE)
	for {
		udp_conn.SetDeadline(time.Now().Add(time.Second * time.Duration(CONN_TIMEOUT)))
		n, err := udp_conn.Read(buff)
		if err != nil {
			log.Println("Connect:listenUDP:Read:", err)
			return
		}
		this.Recv(buff[:n])
	}
}

//Receive udp data
func (this *Connect) Recv(pkg []byte) {
	if len(pkg) < 5 {
		return
	}
	conn_id := uint16(pkg[0])<<8 + uint16(pkg[1])
	sig := uint8(pkg[4])
	switch {
	case sig == SIGNAL_RETRY:
		if udpp, exist := this.out[conn_id]; exist {
			if udpp.Sig != nil {
				udpp.Sig <- [2]uint8{SIGNAL_RETRY, pkg[2]}
			}
		}
	case sig == SIGNAL_DONE:
		if udpp, exist := this.out[conn_id]; exist {
			if udpp.Sig != nil {
				udpp.Sig <- [2]uint8{SIGNAL_DONE}
			}
		}
	case sig == SIGNAL_SEND:
		if udpp, exist := this.in[conn_id]; exist {
			udpp.Recv(pkg)
		} else {
			this.newUDPPackage(false, conn_id).Recv(pkg)
		}
	default:
		log.Println("Connect:Recv:signal error")
	}
}
