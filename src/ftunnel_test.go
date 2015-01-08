package ftunnel

import (
	"bytes"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

type tunnel_test struct {
	request  []byte
	response []byte
}

var (
	c_tcp_addr     = "127.0.0.1:12345"
	s_tcp_addr     = "127.0.0.1:34567"
	s_udp_addr_str = "127.0.0.1:56789"
	s_udp_addr     *net.UDPAddr
	request_nums   = 20
	wg             sync.WaitGroup
	test_data      = []tunnel_test{
		{[]byte("request begin"), []byte("response begin")},
		{[]byte("request continue"), []byte("response continue")},
		{[]byte("request end"), []byte("response end")},
	}
)

//server
func serverRun() {
	udp_listener, err := net.ListenUDP("udp", s_udp_addr)
	if err != nil {
		log.Fatal("serverRun:Listen:", err)
	}
	buff := make([]byte, BUFFER_MAXSIZE)
	//connMap := make(map[string]*Connect)
	for {
		n, u_addr, err := udp_listener.ReadFromUDP(buff)
		if err != nil {
			log.Fatal("serverRun:ReadFromUDP:", err)
		}
		//u_addr_str := u_addr.String()
		//if c, exist := connMap[u_addr_str]; exist {
		//	c.Recv(buff[:n])
		//} else {
		t_conn, err := net.Dial("tcp", s_tcp_addr)
		if err != nil {
			log.Fatal("serverRun:Dial tcp:", err)
		}

		u_conn, err := net.DialUDP("udp", nil, u_addr)
		if err != nil {
			log.Fatal("serverRun:Dial udp:", err)
		}

		conn, err := NewConnect(t_conn, u_conn)
		if err != nil {
			log.Fatal("serverRun:NewConnect:", err)
		}
		//connMap[u_addr_str] = conn
		go conn.Serve()
		conn.Recv(buff[:n])
		//}
	}
}

//client
func clientRun() {
	ln, err := net.Listen("tcp", c_tcp_addr)
	if err != nil {
		log.Fatal("clientRun:Listen:", err)
	}
	for {
		t_conn, err := ln.Accept()
		if err != nil {
			log.Fatal("clientRun:Accept:", err)
		}

		u_conn, err := net.DialUDP("udp", nil, s_udp_addr)
		if err != nil {
			log.Fatal("clientRun:DialUDP:", err)
		}
		conn, err := NewConnect(t_conn, u_conn)
		if err != nil {
			log.Fatal("clientRun:NewConnect:", err)
		}
		_, err = u_conn.Write([]byte{})
		if err != nil {
			log.Fatal("clientRun:NewConnect:Write test:", err)
		}

		u_addr, err := net.ResolveUDPAddr("udp", u_conn.LocalAddr().String())
		if err != nil {
			log.Fatal("clientRun:ResolveUDPAddr:", err)
		}
		udp_conn, err := net.ListenUDP("udp", u_addr)
		if err != nil {
			log.Fatal("clientRun:ListenUDP:", err)
		}

		go conn.ListenUDP(udp_conn)
		go conn.Serve()
	}
}

//request
func requestRun() {
	conn, err := net.DialTimeout("tcp", c_tcp_addr, 1e9)
	if err != nil {
		log.Fatal("requestRun:Dial:", err)
	}
	for _, td := range test_data {
		err := conn.SetDeadline(time.Now().Add(1e9))
		if err != nil {
			log.Fatal("requestRun:SetDeadline:", err)
		}
		_, err = conn.Write(td.request)
		if err != nil {
			log.Fatal("requestRun:Write:", err)
		}
		buff := make([]byte, BUFFER_MAXSIZE)
		n, err := conn.Read(buff)
		if err != nil {
			log.Fatal("requestRun:Read:", err)
		}
		response := buff[:n]
		if bytes.Compare(response, td.response) != 0 {
			log.Fatal("want:", string(td.response), ", but get:", string(response))
		}
	}
	wg.Done()
}

//handle s_tcp_addr, like squid
func handleConn(conn net.Conn) {
	buff := make([]byte, BUFFER_MAXSIZE)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Fatal("Test:Read:", err)
		}
		request := buff[:n]
		notFound := true
		for _, td := range test_data {
			if bytes.Compare(request, td.request) == 0 {
				_, err := conn.Write(td.response)
				if err != nil {
					log.Fatal("Test:Write:", err)
				}
				notFound = false
			}
		}
		if notFound {
			log.Fatal("Test:notFound:", string(request))
		}
	}
}

//main test
func Test(t *testing.T) {
	var err error
	s_udp_addr, err = net.ResolveUDPAddr("udp", s_udp_addr_str)
	if err != nil {
		t.Fatal("Test:ResolveUDPAddr:", err)
	}

	go serverRun()
	go clientRun()

	//make some request
	for i := 0; i < request_nums; i++ {
		go requestRun()
		wg.Add(1)
	}

	go func() {
		//listen s_tcp_addr
		ln, err := net.Listen("tcp", s_tcp_addr)
		if err != nil {
			t.Fatal("Test:Listen:", err)
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				t.Fatal("Test:Accept:", err)
			}
			go handleConn(conn)
		}
	}()

	wg.Wait()
}
