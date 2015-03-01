// generate node identity
package ftunnel

import (
	"hash/fnv"
	"log"
	"net"
	"os"
)

var (
	_nodeId uint64 = NodeIdentity()
)

func NodeIdentity() (r uint64) {
	r = uint64(os.Getpid())
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("E(NodeIdentity):", err)
		return
	}
	for _, inter := range interfaces {
		h := fnv.New32()
		h.Write([]byte(inter.HardwareAddr))
		checksum := h.Sum32()

		r = uint64(checksum)<<32 | r
		return
	}
	return
}
