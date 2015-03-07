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
