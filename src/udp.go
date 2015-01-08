package ftunnel

import (
	"errors"
	"log"
	"net"
	"time"
)

var (
	SplitSize  = 1024
	Timeout    = 10       //seconds signal wait
	CheckPoint = uint8(1) //check and request retry when the remaining pkg numbers less than this number
)

const (
	SIGNAL_SEND  = uint8(0)
	SIGNAL_DONE  = uint8(1)
	SIGNAL_RETRY = uint8(2)
)

//extra data add to udp package
type extraHeader struct {
	conn_id    uint16 //byte 0,1
	pkg_id     uint8  //byte 2
	pkg_max_id uint8  //byte 3
	signal     uint8  //byte 4
}

//NewExtraHeader
func NewExtraHeader(conn_id uint16, pkg_id, pkg_max_id, signal uint8) *extraHeader {
	return &extraHeader{conn_id: conn_id, pkg_id: pkg_id, pkg_max_id: pkg_max_id, signal: signal}
}

//attach header to data
func (this *extraHeader) attach(data []byte) []byte {
	var tmp []byte
	tmp = append(tmp, byte(this.conn_id>>8))
	tmp = append(tmp, byte(this.conn_id))
	tmp = append(tmp, byte(this.pkg_id))
	tmp = append(tmp, byte(this.pkg_max_id))
	tmp = append(tmp, byte(this.signal))
	return append(tmp, data...)
}

//UDP Package
type UDPPackage struct {
	Sig      chan [2]uint8
	conn_tcp net.Conn
	conn_udp *net.UDPConn
	conn_id  uint16
	pkg_buff map[uint8][]byte
	pkg_len  uint8
}

//NewUDPPackage
func NewUDPPackage(conn_tcp net.Conn, conn_udp *net.UDPConn, conn_id uint16) *UDPPackage {
	return &UDPPackage{conn_tcp: conn_tcp, conn_udp: conn_udp, conn_id: conn_id}
}

//send data by udp
func (this *UDPPackage) Send(buff []byte) (err error) {
	//split package
	var pkg_max_id uint8
	pkg_map := make(map[uint8][]byte)
	overflow_flag := true

	for pkg_id := uint8(0); pkg_id < 255; pkg_id++ {
		if len(buff) > SplitSize {
			tmp := buff[:SplitSize]
			buff = buff[SplitSize:]
			pkg_map[pkg_id] = tmp
			continue
		}
		pkg_map[pkg_id] = buff
		overflow_flag = false
		pkg_max_id = pkg_id
		break
	}

	if overflow_flag {
		err = errors.New("UDPPackage:Send:package overflow max size")
		return
	}

	//send package
	for pkg_id, data := range pkg_map {
		header := NewExtraHeader(this.conn_id, pkg_id, pkg_max_id, SIGNAL_SEND)
		err = this.sendUDP(data, header)
		if err != nil {
			return
		}
	}

	//keep pkgs, wait signal
	this.Sig = make(chan [2]uint8)
	go this.handleSignal(pkg_max_id, pkg_map)
	return
}

//attach extra header and send
func (this *UDPPackage) sendUDP(data []byte, header *extraHeader) (err error) {
	_, err = this.conn_udp.Write(header.attach(data))
	return
}

//wait signal
func (this *UDPPackage) handleSignal(pkg_max_id uint8, pkg_map map[uint8][]byte) {
	for {
		select {
		case s := <-this.Sig:
			switch {
			case s[0] == SIGNAL_RETRY:
				header := NewExtraHeader(this.conn_id, s[1], pkg_max_id, SIGNAL_SEND)
				err := this.sendUDP(pkg_map[s[1]], header)
				if err != nil {
					log.Println("UDPPackage:handleSignal:SIGNAL_RETRY:", err)
				}
				continue
			case s[0] == SIGNAL_DONE:
				//do something
			default:
				log.Println("UDPPackage:handleSignal:unknow signal.")
			}
		case <-time.After(time.Duration(Timeout) * time.Second):
			log.Println("UDPPackage:handleSignal:timeout:", this.conn_id)
		}
		this.Sig = nil
		break
	}
}

//receive data
func (this *UDPPackage) Recv(pkg []byte) {
	pkg_id := pkg[2]
	pkg_max_id := uint8(pkg[3])
	if this.pkg_buff == nil {
		this.pkg_buff = make(map[uint8][]byte)
		this.pkg_len = pkg_max_id + 1
		//init timer
		time.AfterFunc(time.Duration(Timeout)*time.Second, this.cleanBuff)
	}

	if _, exist := this.pkg_buff[pkg_id]; !exist {
		this.pkg_buff[pkg_id] = pkg[5:]
		this.pkg_len = this.pkg_len - 1
		if this.pkg_len == 0 {
			//send OK
			header := NewExtraHeader(this.conn_id, 0, 0, SIGNAL_DONE)
			err := this.sendUDP([]byte{}, header)
			if err != nil {
				log.Println("UDPPackage:Recv:SIGNAL_DONE:", err)
			}
			//send to tcp
			buff := []byte{}
			for i := uint8(0); i <= pkg_max_id; i++ {
				if p, exist := this.pkg_buff[i]; exist {
					buff = append(buff, p...)
				} else {
					log.Println("UDPPackage:Recv:BUG!!!")
				}
			}
			_, err = this.conn_tcp.Write(buff)
			if err != nil {
				log.Println("UDPPackage:Recv:write:", err)
			}
			this.cleanBuff()
		}
	}
	if this.pkg_len < CheckPoint {
		for i := uint8(0); i <= pkg_max_id; i++ {
			if _, exist := this.pkg_buff[i]; !exist {
				header := NewExtraHeader(this.conn_id, i, pkg_max_id, SIGNAL_RETRY)
				err := this.sendUDP([]byte{}, header)
				if err != nil {
					log.Println("UDPPackage:Recv:SIGNAL_RETRY:", err)
				}
			}
		}
	}
}

//clean buff
func (this *UDPPackage) cleanBuff() {
	this.pkg_len = 0
	this.pkg_buff = nil
}
