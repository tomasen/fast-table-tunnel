// Handle data transfer
package ftunnel

import (
	"encoding/binary"
	flatbuffers "github.com/google/flatbuffers/go"
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
	CMD_QUERY_IDENTITY uint16 = iota
	CMD_ANSWER_IDENTITY
	CMD_PING
	CMD_PONG
	CMD_CONN
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
	for {
		var b []byte
		read, err := tr.Read(b)
		if err != nil {
			tr.Close()
			log.Println("N(core.ReadNextPacket):", err)
			return nil
		}
		tr.readBytes += read
		tr.readBuffer = append(tr.readBuffer, b...)
		if tr.readBytes > 0 {
			packetLen, packetStart := binary.Uvarint(tr.readBuffer)
			packetSize := int(packetLen) + packetStart
			if packetStart > 0 && tr.readBytes >= packetSize {
				// unpack
				pack := GetRootAsPacket(tr.readBuffer[packetStart:packetLen], 0)

				tr.readBytes -= packetSize
				tr.readBuffer = tr.readBuffer[packetSize:]

				return pack
			}
			// keep reading
		}
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
		if pack == nil {
			break
		}
		builder := flatbuffers.NewBuilder(0)
		PacketStart(builder)
		switch pack.Command() {
		case CMD_PING:
			PacketAddCommand(builder, CMD_PONG)
		case CMD_QUERY_IDENTITY:
			// reply this node's identity
			b := make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(b, _nodeId)

			PacketAddCommand(builder, CMD_ANSWER_IDENTITY)
			PacketAddContentData(builder, b[:n])
		}

		PacketEnd(builder)
		tr.WritePacketBytes(builder.Bytes)
	}
}

func (tr *Transporter) QueryIdentity() uint64 {
	// send QUERY_IDENTITY
	builder := flatbuffers.NewBuilder(0)
	PacketStart(builder)
	PacketAddCommand(builder, CMD_QUERY_IDENTITY)
	PacketEnd(builder)

	tr.WritePacketBytes(builder.Bytes)

	p := tr.ReadNextPacket()
	if p != nil {
		b := p.ContentData()
		identity, _ := binary.Uvarint(b)
		return identity
	}

	return 0
}

func (tr *Transporter) Ping() int64 {
	// send CMD_PING
	builder := flatbuffers.NewBuilder(0)
	PacketStart(builder)
	PacketAddCommand(builder, CMD_PING)
	PacketEnd(builder)
	tr.WritePacketBytes(builder.Bytes)
	s := time.Now()
	// reply and record the latency
	p := tr.ReadNextPacket()
	if p != nil {
		return time.Now().Sub(s).Nanoseconds()
	}
	// return proper error
	return math.MaxInt64
}
