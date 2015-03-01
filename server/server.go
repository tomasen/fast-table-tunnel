package main

import (
	"ftunnel"
	"log"
	"net"
	"sync"
	"time"
)

var (
	BUFFER_MAXSIZE = 1048576
	CLEAN_MAP      = 15
	tcp_addr       = "127.0.0.1:55552"
	udp_addr_str   = "0.0.0.0:2468"
	udp_addr       *net.UDPAddr
	connMap        *ConnMap
)

type ConnMap struct {
	locker sync.Mutex
	cm     map[string]*ftunnel.Connect
}

func NewConnMap() *ConnMap {
	return &ConnMap{cm: make(map[string]*ftunnel.Connect)}
}

func (this *ConnMap) Get(u_addr_str string) *ftunnel.Connect {
	this.locker.Lock()
	defer this.locker.Unlock()

	if c, exist := this.cm[u_addr_str]; exist {
		return c
	}
	return nil
}
func (this *ConnMap) Set(u_addr_str string, c *ftunnel.Connect) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.cm[u_addr_str] = c
}
func (this *ConnMap) Clean() {
	this.locker.Lock()
	defer this.locker.Unlock()

	//	la := len(this.cm)
	for u_addr_str, c := range this.cm {
		if c.Closed {
			delete(this.cm, u_addr_str)
		}
	}
	//log.Println("[info]Clean:", len(this.cm), "/", la)
}

func main() {
	var err error
	udp_addr, err = net.ResolveUDPAddr("udp4", udp_addr_str)
	if err != nil {
		log.Fatal("main:ResolveUDPAddr:", err)
	}

	udp_listener, err := net.ListenUDP("udp4", udp_addr)
	if err != nil {
		log.Fatal("main:Listen:", err)
	}

	ticker := time.NewTicker(time.Second * time.Duration(CLEAN_MAP))
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatal("main:ticker:", r)
			}
		}()

		for _ = range ticker.C {
			connMap.Clean()
		}
	}()

	connMap = NewConnMap()
	for {
		buff := make([]byte, BUFFER_MAXSIZE)
		n, u_addr, err := udp_listener.ReadFromUDP(buff)
		if err != nil {
			log.Fatal("main:ReadFromUDP:", err)
		}
		if n < 5 {
			continue
		}
		u_addr_str := u_addr.String()
		if c := connMap.Get(u_addr_str); c != nil {
			go c.Recv(buff[:n])
		} else if buff[4] != ftunnel.SIGNAL_RESEND {
			t_conn, err := net.Dial("tcp4", tcp_addr)
			if err != nil {
				log.Fatal("main:Dial tcp:", err)
			}

			conn, err := ftunnel.NewConnect(t_conn, udp_listener, u_addr)
			if err != nil {
				log.Fatal("main:NewConnect:", err)
			}
			connMap.Set(u_addr_str, conn)
			go conn.Serve()
			go conn.Recv(buff[:n])
		}
	}
}
