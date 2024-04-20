package cluster

import (
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type AgentStatus uint8

const (
	AgentStatus_Idle AgentStatus = iota
	AgentStatus_Register
	AgentStatus_Lease // 续约
)

// 模块的代理
type ModuleAgent struct {
	module  rpc3.Module
	status  AgentStatus
	leaseID clientv3.LeaseID
}

// 全局注册到etcd3中
func (a *ModuleAgent) Init(t rpc3.ModuleType) {
	a.module.MType = t
	a.module.ClusterID = GCluster.ClusterInfo.Id()

	go a.RunAgent()
}

func (a *ModuleAgent) RegisterAgent() {
	a.module.ID = (a.module.ID) % GCluster.GetModuleMax(a.module.MType)

	if GCluster.moduleMgr.Register(&a.module, a) {
		// 注册成功
		a.status = AgentStatus_Lease

		entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, fmt.Sprintf("%s.OnModuleRegister", a.module.MType.String()))
		logm.DebugfE("模块注册成功: type:%s, ID:%d, clusterID:%s", a.module.MType.String(), a.module.ID, a.module.String())

		time.Sleep(time.Duration(GCluster.moduleMgr.timeGrant / 3))
	} else if GCluster.IsEnough(a.module.MType) {
		a.status = AgentStatus_Idle
	}
}

func (a *ModuleAgent) Lease() {
	if !GCluster.moduleMgr.Lease(a) {
		a.status = AgentStatus_Idle

		entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, fmt.Sprintf("%s.OnModuleUnRegister", a.module.MType.String()))
		logm.DebugfE("模块续约失败: %s %d %s", a.module.MType.String(), a.module.ID, a.module.String())
	} else {
		// 避免cpu忙
		time.Sleep(time.Duration(GCluster.moduleMgr.timeGrant / 3))
	}
}

func (a *ModuleAgent) Idle() {
	// 已经加载足够的模块了
	if GCluster.IsEnough(a.module.MType) {
		return
	}

	a.status = AgentStatus_Register
}

func (a *ModuleAgent) RunAgent() {
	for {
		switch a.status {
		case AgentStatus_Idle:
			a.Idle()
		case AgentStatus_Register:
			a.RegisterAgent()
		case AgentStatus_Lease:
			a.Lease()
		}

		// 暂停100毫秒
		time.Sleep(time.Millisecond * 100)
	}
}
