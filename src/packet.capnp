@0xca8d3dea8c5c8e21

struct Packet {
	version     @0 :UInt16;       # 节点版本号
	conn_id     @0 :Uint64;       # 连接ID，由监听角色（Listener）节点根据本机信息和时间生成
	packet_id   @0 :UInt16;       # 数据包ID，递增
	dst_node    @0 :Uint64;       # 数据包目的地节点：监听角色（Listener）节点或指定的连接角色（Connector）节点，或节点组别。
	src_node    @0 :Uint64;       # 数据包发起者节点ID
	dst_address @0 :Text;         # 目标IP:Port地址
	dst_network @0 :Text;         # 目标网络
	route_list  @0 :List(Uint64); # 已经经过的节点列表（包括本节点）
	command     @0 :UInt16;       # 包指令
	properties  @0 :UInt16;       # 数据包属性：是否加密，是上行包还是下行包，是否压缩
	content     @0 :Data;         # 数据包数据内容
}

