// +build ignore
// extend flatbuffer Packet
package ftunnel

import (
	"bytes"
	capn "github.com/glycerine/go-capnproto"
	"log"
)

func BuildConnPacket(network string, address string) []byte {
	s := capn.NewBuffer(nil)
	d := NewRootPacket(s)
	d.SetCommand(CMD_CONN)
	d.SetDstAddress(address)
	d.SetDstNetwork(network)
	buf := bytes.Buffer{}
	s.WriteToPacked(&buf)
	return buf.Bytes()
}

func BuildDataPacket(b []byte, nodeid uint64) []byte {
	s := capn.NewBuffer(nil)
	d := NewRootPacket(s)
	d.SetCommand(CMD_DATA)
	d.SetDstNode(nodeid)
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
