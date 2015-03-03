// Handle data transfer
package ftunnel

import (
	"encoding/binary"
	"log"
	"net"
	"sync/atomic"
)

var (
	_connid uint64 = 0
)

const (
	CMD_QUERY_IDENTITY = 1 << iota
	CMD_PING           = 1 << iota
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
			if packetStart > 0 && bytesRead >= (int(packetLen)+packetStart) {
				// TODO: unpack
				pack := GetRootAsPacket(buffer[packetStart:packetLen], 0)

				p := pack.Command()
				if (p & CMD_QUERY_IDENTITY) != 0 {
					// TODO: reply this node's identity
				}
				bytesRead = 0
				buffer = buffer[:0]
			} else {
				// TODO: keep reading
			}
		}
	}
}
