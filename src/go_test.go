// test transporter
package ftunnel

import (
	"errors"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"net"
	"testing"
	"time"
)

const (
	TEST_PORT = "63137"
)

// TODO: test transporter ReadNextPacket
func TestReadNextPacket(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:"+TEST_PORT)
	if err != nil {
		t.Fatal(err)
	}
	exit := make(chan bool, 1)
	go func() {
		defer func() {
			exit <- true
		}()

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
				t.Fatal(errors.New("Command mismatch"))
			}

			pack = tr.ReadNextPacket()
			if pack == nil {
				t.Fatal(errors.New("ReadNextPacket return nil (1)"))
			}
			fmt.Println(pack)
		}
	}()

	conn, err := net.Dial("tcp", "localhost:"+TEST_PORT)
	if err != nil {
		t.Fatal(err)
	}

	tr := NewTransporter(conn)
	// Testing CMD_PING
	builder := flatbuffers.NewBuilder(0)
	PacketStart(builder)
	PacketAddCommand(builder, CMD_PING)
	PacketEnd(builder)
	tr.WritePacketBytes(builder.Bytes)

	// Testing unformatted bytes
	// tr.WritePacketBytes([]byte("123"))

	tr.Close()

	select {
	case <-exit:
	case <-time.After(3 * time.Second):
		t.Fatal(errors.New("timed out"))
	}
}

// TODO: test transporter send ping
