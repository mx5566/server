package rpc3

import (
	"fmt"
	"github.com/mx5566/server/base"
	"hash/crc32"
)

func (x *ClusterInfo) Id() uint32 {
	return crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s:%d", x.Ip, x.Port)))
}

func (x *ClusterInfo) GetTopicServerID() string {
	top := fmt.Sprintf("%s%s/%d", base.ServiceName, x.GetServiceType().String(), x.Id())
	return top
}

func (x *ClusterInfo) GetTopicServerType() string {
	top := fmt.Sprintf("%s%s", base.ServiceName, x.GetServiceType().String())
	return top
}
