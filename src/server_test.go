package ftunnel

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

const (
	TCPORT = "64391"
	TSPORT = "64388"
	SPORT  = "64389"
)

func TestSpeed(b *testing.T) {

	// start tunnel client
	lc, err := net.Listen("tcp", "127.0.0.1:"+TCPORT)
	if err != nil {
		b.Fatal(err)
	}
	cl := make(chan net.Listener, 1)
	cl <- lc
	go HandleListners(cl, "127.0.0.1:"+TSPORT)

	// start tunnel server
	ls, err := net.Listen("tcp", "127.0.0.1:"+TSPORT)
	if err != nil {
		b.Fatal(err)
	}
	cs := make(chan net.Listener, 1)
	cs <- ls
	go HandleListners(cs, "127.0.0.1:"+SPORT)

	// start speed test server
	l, err := net.Listen("tcp", "127.0.0.1:"+SPORT)
	if err != nil {
		b.Fatal(err)
	}
	go func() {
		defer l.Close()
		for {
			// Wait for a connection.
			conn, err := l.Accept()
			if err != nil {
				b.Fatal(err)
			}
			// Handle the connection in a new goroutine.
			// The loop then returns to accepting, so that
			// multiple connections may be served concurrently.
			go func(c net.Conn) {
				// Echo all incoming data.
				io.Copy(c, c)
				// Shut down the connection.
				c.Close()
			}(conn)
		}
	}()
	tSpeed(TCPORT, b)
	tSpeed(SPORT, b)
}

func tSpeed(port string, b *testing.T) float64 {
	t0 := time.Now()
	// begin benchmark
	d, err := net.Dial("tcp", "127.0.0.1:"+TCPORT)
	if err != nil {
		// handle error
		b.Fatal(err)
	}

	var total_read, total_write int = 0, 0

	buf := make([]byte, 1024*1024)
	_, err = rand.Read(buf)
	if err != nil {
		b.Fatal(err)
	}
	//t1 := time.Now()
	go func() {
		for i := 0; i < 100; i++ {
			n, _ := d.Write(buf)
			total_write += n
		}
	}()
	for total_read < total_write || total_read == 0 {
		n, e := d.Read(buf)
		if e != nil {
			//fmt.Println(e)
			break
		}
		total_read += n
	}
	t2 := time.Now()
	spd := float64(total_write) / t2.Sub(t0).Seconds() / 1024 / 1024
	fmt.Printf("%.2fMBPS\n", spd)
	return spd
}
