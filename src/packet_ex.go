// +build ignore
// extend flatbuffer Packet
package ftunnel

import (
	capn "github.com/glycerine/go-capnproto"
	"bytes"
	"log"
)

func BuildConnPacket(network string, address string) []byte {
	// send CMD_CONN
	s := capn.NewBuffer(nil)
	d := NewRootPacket(s)
	d.SetCommand(CMD_CONN)	
	d.SetDstAddress(address)
	d.SetDstNetwork(network)
	buf := bytes.Buffer{}
	s.WriteToPacked(&buf)
	return buf.Bytes()
}

func GetAsPacket(b []byte) (pack *Packet, err error) {
  defer func() {
    if r := recover(); r != nil {
      err, _ = r.(error)
    }
  }()

	buf := bytes.NewBuffer(b)

	s, err := capn.ReadFromPackedStream(buf, nil)
	if err != nil {
		log.Println("E(packet_ex.GetAsPacket.ReadFromPackedStream):", err)
		return
	}
	p := ReadRootPacket(s)
	return &p, nil
}

func BuildCommandPacket(cmd uint16) []byte {
	s := capn.NewBuffer(nil)
	d := NewRootPacket(s)
	d.SetCommand(cmd)	
	buf := bytes.Buffer{}
	s.WriteToPacked(&buf)
	return buf.Bytes()
}
