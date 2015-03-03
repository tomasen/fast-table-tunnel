// automatically generated, do not modify

package ftunnel

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Packet struct {
	_tab flatbuffers.Table
}

func GetRootAsPacket(buf []byte, offset flatbuffers.UOffsetT) *Packet {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Packet{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Packet) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Packet) Version() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Packet) ConnId() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Packet) PacketId() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Packet) DestNode() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Packet) SrcNode() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Packet) RouteList(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j * 8))
	}
	return 0
}

func (rcv *Packet) RouteListLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Packet) Command() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Packet) Properties() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Packet) Content(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j * 1))
	}
	return 0
}

func (rcv *Packet) ContentLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func PacketStart(builder *flatbuffers.Builder) { builder.StartObject(9) }
func PacketAddVersion(builder *flatbuffers.Builder, version uint16) { builder.PrependUint16Slot(0, version, 0) }
func PacketAddConnId(builder *flatbuffers.Builder, connId uint64) { builder.PrependUint64Slot(1, connId, 0) }
func PacketAddPacketId(builder *flatbuffers.Builder, packetId uint16) { builder.PrependUint16Slot(2, packetId, 0) }
func PacketAddDestNode(builder *flatbuffers.Builder, destNode uint64) { builder.PrependUint64Slot(3, destNode, 0) }
func PacketAddSrcNode(builder *flatbuffers.Builder, srcNode uint64) { builder.PrependUint64Slot(4, srcNode, 0) }
func PacketAddRouteList(builder *flatbuffers.Builder, routeList flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(routeList), 0) }
func PacketStartRouteListVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(8, numElems, 8)
}
func PacketAddCommand(builder *flatbuffers.Builder, command uint16) { builder.PrependUint16Slot(6, command, 0) }
func PacketAddProperties(builder *flatbuffers.Builder, properties uint16) { builder.PrependUint16Slot(7, properties, 0) }
func PacketAddContent(builder *flatbuffers.Builder, content flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(content), 0) }
func PacketStartContentVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(1, numElems, 1)
}
func PacketEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
