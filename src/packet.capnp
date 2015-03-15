@0xa979b6ae6132202a;

using Go = import "go.capnp";

$Go.package("ftunnel");

struct CapnPacket {
	version     @0 :UInt16;       # 节点版本号
	connId     @1 :UInt64;       # 连接ID，由监听角色（Listener）节点根据本机信息和时间生成
	packetId   @2 :UInt16;       # 数据包ID，递增
	dstNode    @3 :UInt64;       # 数据包目的地节点：监听角色（Listener）节点或指定的连接角色（Connector）节点，或节点组别。
	srcNode    @4 :UInt64;       # 数据包发起者节点ID
	dstAddress @5 :Text;         # 目标IP:Port地址
	dstNetwork @6 :Text;         # 目标网络
	routeList  @7 :List(UInt64); # 已经经过的节点列表（包括本节点）
	command     @8 :UInt16;       # 包指令
	properties  @9 :UInt16;       # 数据包属性：是否加密，是上行包还是下行包，是否压缩
	content     @10 :Data;         # 数据包数据内容
}

