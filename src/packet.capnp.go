package ftunnel

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
)

type Packet C.Struct

func NewPacket(s *C.Segment) Packet          { return Packet(s.NewStruct(32, 4)) }
func NewRootPacket(s *C.Segment) Packet      { return Packet(s.NewRootStruct(32, 4)) }
func AutoNewPacket(s *C.Segment) Packet      { return Packet(s.NewStructAR(32, 4)) }
func ReadRootPacket(s *C.Segment) Packet     { return Packet(s.Root(0).ToStruct()) }
func (s Packet) Version() uint16             { return C.Struct(s).Get16(0) }
func (s Packet) SetVersion(v uint16)         { C.Struct(s).Set16(0, v) }
func (s Packet) ConnId() uint64              { return C.Struct(s).Get64(8) }
func (s Packet) SetConnId(v uint64)          { C.Struct(s).Set64(8, v) }
func (s Packet) PacketId() uint16            { return C.Struct(s).Get16(2) }
func (s Packet) SetPacketId(v uint16)        { C.Struct(s).Set16(2, v) }
func (s Packet) DstNode() uint64             { return C.Struct(s).Get64(16) }
func (s Packet) SetDstNode(v uint64)         { C.Struct(s).Set64(16, v) }
func (s Packet) SrcNode() uint64             { return C.Struct(s).Get64(24) }
func (s Packet) SetSrcNode(v uint64)         { C.Struct(s).Set64(24, v) }
func (s Packet) DstAddress() string          { return C.Struct(s).GetObject(0).ToText() }
func (s Packet) SetDstAddress(v string)      { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Packet) DstNetwork() string          { return C.Struct(s).GetObject(1).ToText() }
func (s Packet) SetDstNetwork(v string)      { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s Packet) RouteList() C.UInt64List     { return C.UInt64List(C.Struct(s).GetObject(2)) }
func (s Packet) SetRouteList(v C.UInt64List) { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Packet) Command() uint16             { return C.Struct(s).Get16(4) }
func (s Packet) SetCommand(v uint16)         { C.Struct(s).Set16(4, v) }
func (s Packet) Properties() uint16          { return C.Struct(s).Get16(6) }
func (s Packet) SetProperties(v uint16)      { C.Struct(s).Set16(6, v) }
func (s Packet) Content() []byte             { return C.Struct(s).GetObject(3).ToData() }
func (s Packet) SetContent(v []byte)         { C.Struct(s).SetObject(3, s.Segment.NewData(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Packet) MarshalJSON() (bs []byte, err error) { return }

type Packet_List C.PointerList

func NewPacketList(s *C.Segment, sz int) Packet_List {
	return Packet_List(s.NewCompositeList(32, 4, sz))
}
func (s Packet_List) Len() int        { return C.PointerList(s).Len() }
func (s Packet_List) At(i int) Packet { return Packet(C.PointerList(s).At(i).ToStruct()) }
func (s Packet_List) ToArray() []Packet {
	n := s.Len()
	a := make([]Packet, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s Packet_List) Set(i int, item Packet) { C.PointerList(s).Set(i, C.Object(item)) }
