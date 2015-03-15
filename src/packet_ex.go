// extend flatbuffer Packet
package ftunnel

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

func PacketAddContentData(builder *flatbuffers.Builder, b []byte) {
	l := len(b)
	PacketStartContentVector(builder, l)
	for i := 0; i < l; i++ {
		builder.PrependByte(b[i])
	}
	o := builder.EndVector(l)
	PacketAddContent(builder, o)
}

func (rcv *Packet) ContentData() []byte {
	l := rcv.ContentLength()
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = rcv.Content(i)
	}
	return b
}

func BuildConnPacket(network string, address string) []byte {
	// send CMD_CONN
	builder := flatbuffers.NewBuilder(0)
	PacketStart(builder)
	PacketAddCommand(builder, CMD_CONN)
	PacketAddDstNetwork(builder, builder.CreateString(network))
	PacketAddDstAddress(builder, builder.CreateString(address))
	builder.Finish(PacketEnd(builder))
	return builder.Bytes[builder.Head():]
}

func GetAsPacket(buf []byte, offset flatbuffers.UOffsetT) (pack *Packet, err error) {
  defer func() {
    if r := recover(); r != nil {
      err, _ = r.(error)
    }
  }()
	return GetRootAsPacket(buf, offset), nil
}

func BuildCommandPacket(cmd uint16) []byte {
	builder := flatbuffers.NewBuilder(0)
	PacketStart(builder)
	PacketAddCommand(builder, cmd)
	builder.Finish(PacketEnd(builder))
	return builder.Bytes[builder.Head():]
}
