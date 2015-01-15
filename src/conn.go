package ftunnel

import (
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	BUFFER_MAXSIZE = 4096
	CONN_TIMEOUT   = 15 //seconds tcp and udp connect timeout
)

//TCP and UDP connect
type Connect struct {
	in           map[uint16]*UDPPackage
	out          map[uint16]*UDPPackage
	tcp_conn     net.Conn
	udp_conn     *net.UDPConn
	dst_udp_addr *net.UDPAddr
	cur          uint16
	locker       sync.Mutex
}

//NewConnect
func NewConnect(tcp_conn net.Conn, udp_conn *net.UDPConn, dst_udp_addr *net.UDPAddr) (*Connect, error) {
	var err error
	c := &Connect{tcp_conn: tcp_conn, udp_conn: udp_conn, dst_udp_addr: dst_udp_addr, cur: 0}
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
				udpp := NewUDPPackage(this.cur, true)
				this.out[this.cur] = udpp
				return udpp
			}
		}
	} else {
		udpp := NewUDPPackage(conn_id, false)
		this.in[conn_id] = udpp
		return udpp
	}
	return nil
}

//send pkgs
func (this *Connect) sendUDP(pkgs [][]byte) {
	for _, pkg := range pkgs {
		_, err := this.udp_conn.WriteToUDP(pkg, this.dst_udp_addr)
		if err != nil {
			log.Println("Connect:sendUDP:", err)
			return
		}
	}
}

//send to tcp
func (this *Connect) sendTCP(buff []byte) {
	if this.tcp_conn != nil {
		_, err := this.tcp_conn.Write(buff)
		if err != nil {
			log.Println("Connect:sendTCP:", err)
			this.close()
		}
	}
}

//close
func (this *Connect) close() {
	if this.tcp_conn != nil {
		this.tcp_conn.Close()
	}
	this.tcp_conn = nil
}

//Serve
func (this *Connect) Serve() {
	defer this.close()

	for {
		this.tcp_conn.SetDeadline(time.Now().Add(time.Second * time.Duration(CONN_TIMEOUT)))
		buff := make([]byte, BUFFER_MAXSIZE)
		n, err := this.tcp_conn.Read(buff)
		if err != nil {
			eMsg := err.Error()
			if eMsg == "EOF" {
				pkgs, err := this.newUDPPackage(true, 0).Split([]byte(eMsg))
				if err != nil {
					log.Println("Connect:Serve:Split:", err)
					return
				}
				this.sendUDP(pkgs)
				return
			}
			if !strings.Contains(eMsg, "timeout") {
				log.Println("Connect:Serve:Read:", err)
			}
			return
		}

		pkgs, err := this.newUDPPackage(true, 0).Split(buff[:n])
		if err != nil {
			log.Println("Connect:Serve:Split:", err)
			return
		}
		this.sendUDP(pkgs)
	}
}

//ListenUDP and Serve
func (this *Connect) ListenAndServe() {
	go this.Serve()

	for {
		buff := make([]byte, BUFFER_MAXSIZE)
		this.udp_conn.SetDeadline(time.Now().Add(time.Second * time.Duration(CONN_TIMEOUT)))
		n, err := this.udp_conn.Read(buff)
		if err != nil {
			if !strings.Contains(err.Error(), "timeout") {
				log.Println("Connect:listenUDP:Read:", err)
			}
			return
		}

		go this.Recv(buff[:n])
	}
}

//Receive udp data
func (this *Connect) Recv(pkg []byte) {
	if len(pkg) < 5 {
		return
	}

	conn_id := uint16(pkg[0])<<8 + uint16(pkg[1])
	sig := pkg[4]
	switch {
	case sig == SIGNAL_RETRY:
		if udpp, exist := this.out[conn_id]; exist {
			pkgs, err := udpp.Resend()
			if err != nil {
				log.Println("Connect:Recv:SIGNAL_RETRY:", err)
			}
			this.sendUDP(pkgs)
		}
	case sig == SIGNAL_DONE:
		if _, exist := this.out[conn_id]; exist {
			//SIGNAL_DONE maybe not need
		}
	case sig == SIGNAL_SEND:
		var udpp *UDPPackage
		if in_udpp, exist := this.in[conn_id]; exist {
			udpp = in_udpp
		} else {
			udpp = this.newUDPPackage(false, conn_id)
		}
		buff, isAll := udpp.Recv(pkg)
		if isAll {
			this.sendTCP(buff)
		}
	case sig == SIGNAL_RESEND:
		if udpp, exist := this.in[conn_id]; exist {
			buff, isAll := udpp.Recv(pkg)
			if isAll {
				this.sendTCP(buff)
			}
		} else {
			//drop pkg
		}
	default:
		log.Println("Connect:Recv:signal error")
	}
}
