package main

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
	go Serve("127.0.0.1:"+TCPORT, "127.0.0.1:"+TSPORT)

	// start tunnel server
	go Serve("127.0.0.1:"+TSPORT, "127.0.0.1:"+SPORT)

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
	const BUFLEN = 1024*1024
	t0 := time.Now()
	// begin benchmark
	d, err := net.Dial("tcp", "127.0.0.1:"+TCPORT)
	if err != nil {
		// handle error
		b.Fatal(err)
	}

	var total_read, total_write int = 0, 0

	buf := make([]byte, BUFLEN)
	_, err = rand.Read(buf)
	if err != nil {
		b.Fatal(err)
	}
	//t1 := time.Now()
	go func() {
		for i := 0; i < 2; i++ {
			n, _ := d.Write(buf)
			total_write += n
		}
	}()
	dst := make([]byte, BUFLEN)
	for total_read < total_write || total_read == 0 {
		n, e := d.Read(dst)
		if e != nil {
			//fmt.Println(e)
			break
		}
		start := total_read % BUFLEN
		total_read += n
		end := total_read % BUFLEN
		if end > start && string(buf[start:end]) != string(dst[:n]) {
			fmt.Println("not equal", len(buf), n, buf[start:start+22], dst[:22])
			b.Failed()
		}
	}
	t2 := time.Now()
	spd := float64(total_write) / t2.Sub(t0).Seconds() / 1024 / 1024
	fmt.Printf("%.2fMBPS\n", spd)
	return spd
}
