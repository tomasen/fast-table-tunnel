fast-table-tunnel 是成对工作的TCP隧道工具。特点是高速、系统开销小的加密方式

[![Build Status](https://travis-ci.org/tomasen/fast-table-tunnel.svg?branch=master)](https://travis-ci.org/tomasen/fast-table-tunnel)

fast-table-tunnel 支持 Gracefully shutdown

# USAGE #

fast-table-tunnel -c <ip:port> -s <ip:port> -id <service name>

-c 连接目标的IP和端口
-s 本地监听的IP和端口，可以只有端口。例如 ":8080"
-id 服务名称，当同时运行多个 fast-table-tunnel 时 ，用以区分
-log 指定日志文件的路径

kill -HUP <pid>  不间断服务的启动新的可执行文件并指令旧进程  Gracefully shutdown 

# EXAMPLE #

A主机 IP: 1.1.1.1  

B主机 IP: 2.2.2.2  并在 127.0.0.1:3028 上运行了 squid 服务

在A主机上执行 fast-table-tunnel -s 0.0.0.0:8080 -c 2.2.2.2:60000 -id squid

在B主机上执行 fast-table-tunnel -s 0.0.0.0:6000 -c 1.1.1.1:3028 -id squid

之后访问 1.1.1.1:8080 即可使用B主机的squid服务。同时A-B之间的通讯是加密的
