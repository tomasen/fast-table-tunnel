// test transporter
package ftunnel

import (
	"errors"
	"fmt"
//	flatbuffers "github.com/google/flatbuffers/go"
	"log"
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

			pack = tr.ReadNextPacket()
			if pack == nil {
				t.Fatal(errors.New("ReadNextPacket return nil (1)"))
			}
			if pack.Command() != CMD_PONG {
				log.Println("pack.Command()", pack.Command(), CMD_PONG)
				t.Fatal(errors.New("Command mismatch"))
			}

			pack = tr.ReadNextPacket()
			if pack != nil {
				t.Fatal(errors.New("ReadNextPacket return non-nil (3)"))
			}

			pack = tr.ReadNextPacket()
			if pack == nil {
				t.Fatal(errors.New("ReadNextPacket return nil (2)"))
			}
			if pack.Command() != CMD_QUERY_IDENTITY {
				log.Println("pack.Command()", pack.Command(), CMD_QUERY_IDENTITY)
				t.Fatal(errors.New("Command mismatch"))
			}

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

	tr.Close()

	select {
	case <-exit:
		fmt.Println("TestReadNextPacket done")
	case <-time.After(3 * time.Second):
		t.Fatal(errors.New("timed out"))
	}
}

// TODO: test transporter ping pong
