// test transporter
package ftunnel

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	capn "github.com/glycerine/go-capnproto"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"
)

const (
	TEST_PORT = "63137"
)

// test transporter ReadNextPacket
func TestReadNextPacket(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:"+TEST_PORT)
	if err != nil {
		t.Fatal(err)
	}
	exit := make(chan bool, 1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				t.Fatal(err)
			}

			tr := NewTransporter(conn)
			pack := tr.ReadNextPacket()
			if pack == nil {
				t.Fatal(errors.New("ReadNextPacket return nil (0)"))
			}
			if pack.Command() != CMD_PING {
				log.Println("pack.Command()", pack.Command(), CMD_PING)
				t.Fatal(errors.New("Command mismatch"))
			}
			fmt.Println("PASS 0")

			pack = tr.ReadNextPacket()
			if pack == nil {
				t.Fatal(errors.New("ReadNextPacket return nil (1)"))
			}
			if pack.Command() != CMD_PONG {
				log.Println("pack.Command()", pack.Command(), CMD_PONG)
				t.Fatal(errors.New("Command mismatch"))
			}
			fmt.Println("PASS 1")

			pack = tr.ReadNextPacket()
			if pack != nil {
				t.Fatal(errors.New("ReadNextPacket return non-nil (3)"))
			}
			fmt.Println("PASS 2")

			pack = tr.ReadNextPacket()
			if pack == nil {
				t.Fatal(errors.New("ReadNextPacket return nil (2)"))
			}
			if pack.Command() != CMD_QUERY_IDENTITY {
				log.Println("pack.Command()", pack.Command(), CMD_QUERY_IDENTITY)
				t.Fatal(errors.New("Command mismatch"))
			}
			fmt.Println("PASS 3")

			pack = tr.ReadNextPacket()
			if pack == nil {
				t.Fatal(errors.New("ReadNextPacket return nil (2)"))
			}
			b := pack.Content()
			if len(b) < MTU || fmt.Sprintf("%X", md5.Sum(b)) != pack.DstNetwork() {
				t.Fatal(errors.New("ContentData mismatch"))
			}
			fmt.Println("PASS 4")

			ln.Close()
			exit <- true
			break
		}
	}()

	conn, err := net.Dial("tcp", "localhost:"+TEST_PORT)
	if err != nil {
		t.Fatal(err)
	}

	tr := NewTransporter(conn)
	// Testing CMD_PING
	tr.WritePacketBytes(BuildCommandPacket(CMD_PING))

	tr.WritePacketBytes(BuildCommandPacket(CMD_PONG))

	time.Sleep(time.Second)

	// Testing unformatted bytes
	tr.WritePacketBytes([]byte("12"))

	tr.WritePacketBytes(BuildCommandPacket(CMD_QUERY_IDENTITY))

	tc := make([]byte, 3*MTU)
	for i := range tc {
		tc[i] = byte(rand.Intn(1))
	}

	s := capn.NewBuffer(nil)
	d := NewRootPacket(s)
	d.SetCommand(CMD_CONN)
	d.SetDstNetwork(fmt.Sprintf("%X", md5.Sum(tc)))
	d.SetContent(tc)
	buf := bytes.Buffer{}
	s.WriteToPacked(&buf)

	tr.WritePacketBytes(buf.Bytes())

	tr.Close()

	select {
	case <-exit:
		fmt.Println("TestReadNextPacket done")
	case <-time.After(3 * time.Second):
		t.Fatal(errors.New("timed out"))
	}
}

// TODO: test transporter ping pong
func TestTransporterService(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:"+TEST_PORT)
	if err != nil {
		t.Fatal(err)
	}
	exit := make(chan bool, 1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				t.Fatal(err)
			}

			tr := NewTransporter(conn)
			tr.ServConnection()

			exit <- true
			break
		}
	}()

	conn, err := net.Dial("tcp", "localhost:"+TEST_PORT)
	if err != nil {
		t.Fatal(err)
	}

	tr := NewTransporter(conn)

	tr.WritePacketBytes(BuildCommandPacket(CMD_PING))
	s := time.Now()

	// reply and record the latency
	p := tr.ReadNextPacket()
	if p == nil || p.Command() != CMD_PONG {
		t.Fatal("PING PONG service Failed")
	}
	fmt.Println("ping", time.Now().Sub(s).Nanoseconds(), "ns")

	tr.WritePacketBytes(BuildCommandPacket(CMD_QUERY_IDENTITY))
	p = tr.ReadNextPacket()
	if p == nil || p.Command() != CMD_ANSWER_IDENTITY || len(p.Content()) <= 0 {
		t.Fatal("PING PONG service Failed")
	}

	tr.WritePacketBytes(BuildCommandPacket(CMD_CLOSE))

	select {
	case <-exit:
		fmt.Println("TestService done")
	case <-time.After(3 * time.Second):
		t.Fatal(errors.New("timed out"))
	}
}
