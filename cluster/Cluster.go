package cluster

import (
	"context"
	"github.com/mx5566/server/entity"
	"github.com/mx5566/server/rpc"
)

type ClusterInfo struct {
	Ip   string
	Port uint16
}

type Cluster struct {
	entity.Entity
}

func (c *Cluster) Init() {
	c.Entity.Init()

	entity.RegisterEntity(c)
}

func (c *Cluster) AddClusterNode(ctx context.Context, info *rpc.ClusterInfo) {

}

func (c *Cluster) DelClusterNode(ctx context.Context, info *rpc.ClusterInfo) {

}
