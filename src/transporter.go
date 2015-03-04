// Handle data transfer
package ftunnel

import (
	"encoding/binary"
	flatbuffers "github.com/google/flatbuffers/go"
	"log"
	"net"
	"sync/atomic"
)

var (
	_connid uint64 = 0
)

const (
	CMD_QUERY_IDENTITY  = 1 << iota
	CMD_ANSWER_IDENTITY = 1 << iota
	CMD_PING            = 1 << iota
)

type Transporter struct {
	net.Conn
}

func ConnId() uint64 {
	return atomic.AddUint64(&_connid, 1)
}

func (tr *Transporter) HandleConnection() {
	// TODO: unpack and proper reply
	var buffer []byte
	var bytesRead int = 0
	var packetLen uint64 = 0
	var packetStart int = 0
	for {
		var b []byte
		read, err := tr.Read(b)
		if err != nil {
			tr.Close()
			log.Println("N(core.HandleConnection):", err)
			return
		}
		bytesRead += read
		buffer = append(buffer, b...)
		if bytesRead > 0 {
			packetLen, packetStart = binary.Uvarint(buffer)
			packetSize := int(packetLen) + packetStart
			if packetStart > 0 && bytesRead >= packetSize {
				// unpack
				pack := GetRootAsPacket(buffer[packetStart:packetLen], 0)

				p := pack.Command()
				if (p & CMD_QUERY_IDENTITY) != 0 {
					// reply this node's identity
					builder := flatbuffers.NewBuilder(0)
					PacketAddCommand(builder, CMD_ANSWER_IDENTITY)
					var b0 []byte
					binary.PutUvarint(b0, _nodeId)
					PacketStartContentVector(builder, len(b0))
					for i := len(b0); i >= 0; i-- {
						builder.PrependByte(b0[i])
					}
					builder.EndVector(len(b0))
					sz := PacketEnd(builder)
					binary.PutUvarint(b0, uint64(sz))
					tr.Write(b0)
					tr.Write(builder.Bytes)
				}

				bytesRead -= packetSize
				buffer = buffer[packetSize:]
			}
			// keep reading
		}
	}
}

func (tr *Transporter) QueryIdentity() *Packet {
	// send QUERY_IDENTITY
	builder := flatbuffers.NewBuilder(0)
	PacketAddCommand(builder, CMD_QUERY_IDENTITY)
	PacketEnd(builder)
	tr.Write(builder.Bytes)
	// TODO: read reply Packet
	var b []byte
	_, err := tr.Read(b)
	if err != nil {
		tr.Close()
		log.Println("N(core.QueryIdentity.Read):", err)
		return nil
	}
	
	return nil
}

func (tr *Transporter) Ping() {
	// send CMD_PING
	builder := flatbuffers.NewBuilder(0)
	PacketAddCommand(builder, CMD_PING)
	PacketEnd(builder)
	tr.Write(builder.Bytes)
	// TODO: reply and record the latency
}
