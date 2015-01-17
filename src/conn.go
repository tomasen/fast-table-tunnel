package ftunnel

import (
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	BUFFER_MAXSIZE = 1400
	CONN_TIMEOUT   = 10  //seconds tcp and udp connect timeout
	RETRY_TIMEOUT  = 150 //Millisecond
	WRITE_TIMEOUT  = 500 //Millisecond
	RETRY_NUM      = 10
)

type recv_udpp struct {
	udpp  *UDPPackage
	sig   chan bool
	isAll bool
}

//TCP and UDP connect
type Connect struct {
	in           map[uint16]*recv_udpp
	out          map[uint16]*UDPPackage
	tcp_conn     net.Conn
	udp_conn     *net.UDPConn
	dst_udp_addr *net.UDPAddr
	cur          uint16
	tcp_cur      uint16
	locker       sync.Mutex
	Closed       bool
}

//NewConnect
func NewConnect(tcp_conn net.Conn, udp_conn *net.UDPConn, dst_udp_addr *net.UDPAddr) (*Connect, error) {
	var err error
	c := &Connect{tcp_conn: tcp_conn, udp_conn: udp_conn, dst_udp_addr: dst_udp_addr, cur: 0, tcp_cur: 1}
	c.in = make(map[uint16]*recv_udpp)
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
		sig := make(chan bool)
		go this.sendRetry(conn_id, sig)
		this.in[conn_id] = &recv_udpp{udpp: udpp, sig: sig, isAll: false}
		return udpp
	}
	return nil
}

//send retry pkg
func (this *Connect) sendRetry(conn_id uint16, sig chan bool) {
	defer close(sig)

	for i := 1; i <= RETRY_NUM; i++ {
		select {
		case <-time.After(time.Millisecond * time.Duration(RETRY_TIMEOUT*i)):
			if r_udpp, exist := this.in[conn_id]; exist {
				this.sendUDP([][]byte{r_udpp.udpp.Retry()})
			} else {
				break
			}
		case <-sig:
			break
		}
	}
}

//send pkgs
func (this *Connect) sendUDP(pkgs [][]byte) {
	for _, pkg := range pkgs {
		if this.udp_conn != nil {
			this.udp_conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(WRITE_TIMEOUT)))
			_, err := this.udp_conn.WriteToUDP(pkg, this.dst_udp_addr)
			if err != nil {
				log.Println("Connect:sendUDP:", err)
				return
			}
		}
	}
}

//send to tcp
func (this *Connect) sendTCP() {
	for conn_id, r_udpp := range this.in {
		if conn_id == this.tcp_cur && r_udpp.isAll {
			if len(r_udpp.udpp.Buff) == 0 {
				this.closeTCP()
				return
			} else {
				if this.tcp_conn != nil {
					this.tcp_conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(WRITE_TIMEOUT)))
					_, err := this.tcp_conn.Write(r_udpp.udpp.Buff)
					if err != nil {
						log.Println("Connect:sendTCP:", err)
						this.closeTCP()
						return
					}
					delete(this.in, conn_id)
				}
			}
			this.tcp_cur++
		}
	}
}

//closeTCP
func (this *Connect) closeTCP() {
	if this.tcp_conn != nil {
		this.tcp_conn.Close()
	}
	this.tcp_conn = nil
	this.Closed = true
}

//Close all
func (this *Connect) closeUDP() {
	if this.udp_conn != nil {
		this.udp_conn.Close()
	}
	this.udp_conn = nil
	this.Closed = true
}

//Serve
func (this *Connect) Serve() {
	defer this.closeTCP()

	for {
		if this.tcp_conn == nil {
			break
		}
		this.tcp_conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(CONN_TIMEOUT)))
		buff := make([]byte, BUFFER_MAXSIZE)
		n, err := this.tcp_conn.Read(buff)
		if err != nil {
			eMsg := err.Error()
			if eMsg == "EOF" {
				this.sendUDP([][]byte{this.newUDPPackage(true, 0).ClosePkg()})
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
	defer this.closeUDP()

	go this.Serve()

	for {
		buff := make([]byte, BUFFER_MAXSIZE)
		this.udp_conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(CONN_TIMEOUT)))
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
	this.locker.Lock()
	defer this.locker.Unlock()

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
		} else {
			log.Println("[debug]SIGNAL_RETRY not found:", conn_id)
		}
	case sig == SIGNAL_DONE:
		if _, exist := this.out[conn_id]; exist {
			//SIGNAL_DONE maybe not need
		}
	case sig == SIGNAL_SEND:
		var udpp *UDPPackage
		if r_udpp, exist := this.in[conn_id]; exist {
			udpp = r_udpp.udpp
		} else {
			this.locker.Unlock()
			udpp = this.newUDPPackage(false, conn_id)
			this.locker.Lock()
		}
		if udpp.Recv(pkg) {
			if r_udpp, exist2 := this.in[conn_id]; exist2 {
				r_udpp.isAll = true
				r_udpp.sig <- true
			} else {
				log.Println("Connect:SIGNAL_SEND:BUG!!!")
			}
			this.sendTCP()
		}
	case sig == SIGNAL_RESEND:
		if r_udpp, exist := this.in[conn_id]; exist {
			if r_udpp.udpp.Recv(pkg) {
				r_udpp.isAll = true
				r_udpp.sig <- true
			}
		} else {
			//drop pkg
			//log.Println("[debug]SIGNAL_RESEND drop:", conn_id)
		}
		this.sendTCP()
	case sig == SIGNAL_CLOSE:
		this.locker.Unlock()
		this.newUDPPackage(false, conn_id)
		this.locker.Lock()
		if r_udpp, exist := this.in[conn_id]; exist {
			r_udpp.isAll = true
			r_udpp.sig <- true
		} else {
			log.Println("Connect:SIGNAL_CLOSE:BUG!!!")
		}
		this.sendTCP()
	default:
		log.Println("Connect:Recv:signal error")
	}
}
