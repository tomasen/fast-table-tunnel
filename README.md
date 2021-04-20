fast-table-tunnel 是成对工作的TCP隧道工具。特点是配置简单，混淆方式系统开销小。

# USAGE #

fast-table-tunnel -s <ip:port> -c <ip:port>

-s 本地监听的IP和端口，可以只有端口。例如 ":8080"
-c 连接目标的IP和端口

# EXAMPLE #

例如 A主机 IP: 1.1.1.1，B主机 IP: 2.2.2.2，并在 127.0.0.1:1194 上运行了 openvpn 服务。

在A主机上执行 `fast-table-tunnel -s :1194 -c 2.2.2.2:60000`。在B主机上执行 `fast-table-tunnel -s :6000 -c 127.0.0.1:1194`。 

部署完成后，访问 1.1.1.1:1194 即可使用B主机的 openvpn 服务。同时A-B之间的通讯是加入混淆干扰的。
