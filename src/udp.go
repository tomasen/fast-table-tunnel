package ftunnel

import (
	"errors"
	"log"
	"sync"
)

var (
	SplitSize = 1024
)

const (
	SIGNAL_SEND   = uint8(0)
	SIGNAL_DONE   = uint8(1)
	SIGNAL_RETRY  = uint8(2)
	SIGNAL_RESEND = uint8(3)
	SIGNAL_CLOSE  = uint8(4)
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
	conn_id      uint16
	pkg_in_buff  map[uint8][]byte
	pkg_out_buff map[uint8][]byte
	pkg_max_id   uint8
	recv_lock    sync.Mutex
	Buff         []byte
}

//NewUDPPackage
func NewUDPPackage(conn_id uint16, isOut bool) *UDPPackage {
	udpp := &UDPPackage{conn_id: conn_id, Buff: []byte{}}
	if isOut {
		udpp.pkg_out_buff = make(map[uint8][]byte)
	} else {
		udpp.pkg_in_buff = make(map[uint8][]byte)
	}
	return udpp
}

//Split data and return pkgs
func (this *UDPPackage) Split(buff []byte) (pkgs [][]byte, err error) {
	//split package
	overflow_flag := true

	for pkg_id := uint8(0); pkg_id < 255; pkg_id++ {
		if len(buff) > SplitSize {
			tmp := buff[:SplitSize]
			buff = buff[SplitSize:]
			this.pkg_out_buff[pkg_id] = tmp
			continue
		}
		this.pkg_out_buff[pkg_id] = buff
		overflow_flag = false
		this.pkg_max_id = pkg_id
		break
	}

	if overflow_flag {
		err = errors.New("UDPPackage:Send:package overflow max size")
		return
	}

	//ready pkgs
	for pkg_id, data := range this.pkg_out_buff {
		header := NewExtraHeader(this.conn_id, pkg_id, this.pkg_max_id, SIGNAL_SEND)
		pkgs = append(pkgs, header.attach(data))
		if err != nil {
			return
		}
	}

	return
}

//Get need retry pkg
func (this *UDPPackage) Retry() []byte {
	return NewExtraHeader(this.conn_id, 0, 0, SIGNAL_RETRY).attach([]byte{})
}

//Get resend pkgs
func (this *UDPPackage) Resend() (pkgs [][]byte, err error) {
	for pkg_id, data := range this.pkg_out_buff {
		header := NewExtraHeader(this.conn_id, pkg_id, this.pkg_max_id, SIGNAL_RESEND)
		pkgs = append(pkgs, header.attach(data))
		if err != nil {
			return
		}
	}
	return
}

//receive data
func (this *UDPPackage) Recv(pkg []byte) bool {
	this.recv_lock.Lock()
	defer this.recv_lock.Unlock()

	pkg_id := pkg[2]
	pkg_max_id := uint8(pkg[3])
	if _, exist := this.pkg_in_buff[pkg_id]; !exist {
		this.pkg_in_buff[pkg_id] = pkg[5:]
		if uint8(len(this.pkg_in_buff)-1) == pkg_max_id {
			for i := uint8(0); i <= pkg_max_id; i++ {
				if p, exist := this.pkg_in_buff[i]; exist {
					this.Buff = append(this.Buff, p...)
				} else {
					log.Println("UDPPackage:Recv:BUG!!!")
				}
			}
			return true
		}
	}
	return false
}

//Get close connect pkg
func (this *UDPPackage) ClosePkg() []byte {
	return NewExtraHeader(this.conn_id, 0, 0, SIGNAL_CLOSE).attach([]byte{})
}
