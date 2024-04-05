package rpc3

import (
	"fmt"
	"hash/crc32"
)

func (x *ClusterInfo) Id() uint32 {
	return crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s:%d", x.Ip, x.Port)))
}
