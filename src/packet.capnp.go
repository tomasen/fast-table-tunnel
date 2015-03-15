package ftunnel

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
)

type CapnPacket C.Struct

func NewCapnPacket(s *C.Segment) CapnPacket      { return CapnPacket(s.NewStruct(32, 4)) }
func NewRootCapnPacket(s *C.Segment) CapnPacket  { return CapnPacket(s.NewRootStruct(32, 4)) }
func AutoNewCapnPacket(s *C.Segment) CapnPacket  { return CapnPacket(s.NewStructAR(32, 4)) }
func ReadRootCapnPacket(s *C.Segment) CapnPacket { return CapnPacket(s.Root(0).ToStruct()) }
func (s CapnPacket) Version() uint16             { return C.Struct(s).Get16(0) }
func (s CapnPacket) SetVersion(v uint16)         { C.Struct(s).Set16(0, v) }
func (s CapnPacket) ConnId() uint64              { return C.Struct(s).Get64(8) }
func (s CapnPacket) SetConnId(v uint64)          { C.Struct(s).Set64(8, v) }
func (s CapnPacket) PacketId() uint16            { return C.Struct(s).Get16(2) }
func (s CapnPacket) SetPacketId(v uint16)        { C.Struct(s).Set16(2, v) }
func (s CapnPacket) DstNode() uint64             { return C.Struct(s).Get64(16) }
func (s CapnPacket) SetDstNode(v uint64)         { C.Struct(s).Set64(16, v) }
func (s CapnPacket) SrcNode() uint64             { return C.Struct(s).Get64(24) }
func (s CapnPacket) SetSrcNode(v uint64)         { C.Struct(s).Set64(24, v) }
func (s CapnPacket) DstAddress() string          { return C.Struct(s).GetObject(0).ToText() }
func (s CapnPacket) SetDstAddress(v string)      { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s CapnPacket) DstNetwork() string          { return C.Struct(s).GetObject(1).ToText() }
func (s CapnPacket) SetDstNetwork(v string)      { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s CapnPacket) RouteList() C.UInt64List     { return C.UInt64List(C.Struct(s).GetObject(2)) }
func (s CapnPacket) SetRouteList(v C.UInt64List) { C.Struct(s).SetObject(2, C.Object(v)) }
func (s CapnPacket) Command() uint16             { return C.Struct(s).Get16(4) }
func (s CapnPacket) SetCommand(v uint16)         { C.Struct(s).Set16(4, v) }
func (s CapnPacket) Properties() uint16          { return C.Struct(s).Get16(6) }
func (s CapnPacket) SetProperties(v uint16)      { C.Struct(s).Set16(6, v) }
func (s CapnPacket) Content() []byte             { return C.Struct(s).GetObject(3).ToData() }
func (s CapnPacket) SetContent(v []byte)         { C.Struct(s).SetObject(3, s.Segment.NewData(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s CapnPacket) MarshalJSON() (bs []byte, err error) { return }

type CapnPacket_List C.PointerList

func NewCapnPacketList(s *C.Segment, sz int) CapnPacket_List {
	return CapnPacket_List(s.NewCompositeList(32, 4, sz))
}
func (s CapnPacket_List) Len() int            { return C.PointerList(s).Len() }
func (s CapnPacket_List) At(i int) CapnPacket { return CapnPacket(C.PointerList(s).At(i).ToStruct()) }
func (s CapnPacket_List) ToArray() []CapnPacket {
	n := s.Len()
	a := make([]CapnPacket, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s CapnPacket_List) Set(i int, item CapnPacket) { C.PointerList(s).Set(i, C.Object(item)) }
