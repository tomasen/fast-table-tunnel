// Handle data transfer
package ftunnel

import (
	"bytes"
	"encoding/binary"
	capn "github.com/glycerine/go-capnproto"
	"log"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_connid uint64 = 0
)

const (
	CMD_QUERY_IDENTITY = iota
	CMD_ANSWER_IDENTITY
	CMD_PING
	CMD_PONG
	CMD_CONN
	CMD_DATA
	CMD_CLOSE
)

const (
	MTU = 1500
)

type Transporter struct {
	net.Conn
	m          *sync.Mutex
	readBuffer []byte
	readBytes  int
}

func ConnId() uint64 {
	return atomic.AddUint64(&_connid, 1)
}

func NewTransporter(conn net.Conn) *Transporter {
	return &Transporter{conn, &sync.Mutex{}, []byte{}, 0}
}

func (tr *Transporter) ReadNextPacket() *Packet {
	tr.m.Lock()
	defer tr.m.Unlock()
	b := make([]byte, MTU)
	for {
		if tr.readBytes > 0 {
			packetLen, packetStart := binary.Uvarint(tr.readBuffer)
			packetSize := int(packetLen) + packetStart
			if packetStart > 0 && tr.readBytes >= packetSize {
				// unpack
				// log.Println("N(core.ReadNextPacket.unpacking):", tr.readBytes, packetStart,
				// 	packetLen, packetSize, len(tr.readBuffer))
				pack, err := GetAsPacket(tr.readBuffer[packetStart:packetSize])
				if err != nil {
					log.Println("E(core.ReadNextPacket.GetAsPacket):", err)
				}
				tr.readBytes -= packetSize
				tr.readBuffer = tr.readBuffer[packetSize:]

				return pack
			}
			// keep reading
		}
		//tr.SetDeadline(time.Now().Add(time.Second * 30))
		read, err := tr.Read(b)
		//log.Println("N(core.ReadNextPacket.read):", read)
		if err != nil {
			tr.Close()
			log.Println("N(core.ReadNextPacket):", err)
			return nil
		}
		tr.readBytes += read
		tr.readBuffer = append(tr.readBuffer, b[:read]...)
	}
	return nil
}

func (tr *Transporter) WritePacketBytes(p []byte) {
	tr.m.Lock()
	defer tr.m.Unlock()

	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, uint64(len(p)))
	tr.Write(b[:n])
	tr.Write(p)
}

func (tr *Transporter) ServConnection() {
	for {
		pack := tr.ReadNextPacket()
		if pack == nil || pack.Command() == CMD_CLOSE{
			break
		}
		s := capn.NewBuffer(nil)
		d := NewRootPacket(s)
		switch pack.Command() {
		case CMD_PING:
			d.SetCommand(CMD_PONG)
		case CMD_QUERY_IDENTITY:
			// reply this node's identity
			b := make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(b, _nodeId)

			d.SetCommand(CMD_ANSWER_IDENTITY)
			d.SetContent(b[:n])
		}

		buf := bytes.Buffer{}
		s.WriteToPacked(&buf)
		tr.WritePacketBytes(buf.Bytes())
	}
}

func (tr *Transporter) QueryIdentity() uint64 {
	// send QUERY_IDENTITY
	tr.WritePacketBytes(BuildCommandPacket(CMD_QUERY_IDENTITY))

	p := tr.ReadNextPacket()
	if p != nil {
		b := p.Content()
		identity, _ := binary.Uvarint(b)
		return identity
	}

	return 0
}

func (tr *Transporter) Ping() int64 {
	// send CMD_PING
	tr.WritePacketBytes(BuildCommandPacket(CMD_PING))
	s := time.Now()
	// reply and record the latency
	p := tr.ReadNextPacket()
	if p != nil {
		return time.Now().Sub(s).Nanoseconds()
	}
	// return proper error
	return math.MaxInt64
}
