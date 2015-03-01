// Handle data transfer
package ftunnel

import (
	"sync/atomic"
)

var (
	_connid uint64 = 0
)

func ConnId() uint64 {
	return atomic.AddUint64(&_connid, 1)
}
