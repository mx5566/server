package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/rpc3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

// 管理所有的全局模块 比如账号管理器  邮件管理器  聊天管理器等
type ModuleMgr struct {
	moduleAgents [rpc3.ModuleType_END]map[int64]*rpc3.Module
	moduleMutex  [rpc3.ModuleType_END]sync.Mutex
	client       *clientv3.Client
	lease        clientv3.Lease
	timeGrant    int64
}

func (m *ModuleMgr) Init(endPoints []string, t int64) {
	conf := clientv3.Config{
		Endpoints: endPoints,
		//DialTimeout: time.Duration(t) * time.Second,
	}

	client, err := clientv3.New(conf)
	if err != nil {
		logm.PanicfE("服务注册模块启动失败: %v, 错误: %s\n", endPoints, err.Error())
		return
	}
	lease := clientv3.NewLease(client)

	m.client = client
	m.lease = lease
	m.timeGrant = t

	// 初始化
	for i := 0; i < int(rpc3.ModuleType_END); i++ {
		m.moduleAgents[i] = make(map[int64]*rpc3.Module)
		m.moduleMutex[i] = sync.Mutex{}
	}

	go m.Run()
}

// 模块的发现
func (m *ModuleMgr) Run() {
	watchChan := m.client.Watch(context.Background(), base.ServiceName, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type.String() {
			case "PUT":
				m.AddModule(event.Kv.Value)
			case "DELETE":
				m.DelModule(event.Kv.Value)
			}
		}
	}

	return
}

func (m *ModuleMgr) AddModule(data []byte) {
	module := &rpc3.Module{}
	err := json.Unmarshal(data, module)
	if err != nil {
		logm.ErrorfE("模块的发现数据解析失败:%s", err.Error())
		return
	}

	if module.MType < rpc3.ModuleType_AccountMgr || module.MType >= rpc3.ModuleType_END {
		logm.ErrorfE("模块的发现服务类型错误:%d", module.MType)
		return
	}

	m.moduleMutex[module.MType].Lock()
	m.moduleAgents[module.MType][module.GetID()] = module
	m.moduleMutex[module.MType].Unlock()

	logm.DebugfE("模块服务增加模块成功:%s", module.String())
}

func (m *ModuleMgr) DelModule(data []byte) {
	module := &rpc3.Module{}
	err := json.Unmarshal(data, module)
	if err != nil {
		logm.ErrorfE("模块的发现数据解析失败:%s", err.Error())
		return
	}

	if module.MType < rpc3.ModuleType_AccountMgr || module.MType >= rpc3.ModuleType_END {
		logm.ErrorfE("模块的发现服务类型错误:%d", module.MType)
		return
	}

	m.moduleMutex[module.MType].Lock()
	delete(m.moduleAgents[module.MType], module.GetID())
	m.moduleMutex[module.MType].Unlock()

	logm.DebugfE("模块服务删除模块成功:%s", module.String())

}

func (m *ModuleMgr) GetModuleNum(t rpc3.ModuleType) int {
	m.moduleMutex[t].Lock()
	length := len(m.moduleAgents[t])
	m.moduleMutex[t].Unlock()

	return length
}

func (m *ModuleMgr) Register(module *rpc3.Module, agent *ModuleAgent) bool {
	//设置租约时间
	leaseResp, err := m.lease.Grant(context.Background(), m.timeGrant)
	if err != nil {
		logm.ErrorfE("设置租约时间错误: %s\n", err.Error())
		return false
	}

	agent.leaseID = leaseResp.ID

	key := base.ModuleNameDir + module.MType.String() + "/" + fmt.Sprintf("%d", module.ID)
	val, _ := json.Marshal(module)

	tx := m.client.Txn(context.Background())

	// CAS
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	tx.Then(clientv3.OpPut(key, string(val), clientv3.WithLease(leaseResp.ID)))
	tx.Else()

	_, err = tx.Commit()
	if err != nil {
		logm.ErrorE("etcd3注册模块事务失败: %s %s", err.Error(), module.String())
		return false
	}

	return true
}

func (m *ModuleMgr) Lease(agent *ModuleAgent) bool {
	_, err := m.lease.KeepAliveOnce(context.Background(), agent.leaseID)
	if err != nil {
		logm.ErrorfE("etcd moudle lease KeepAliveOnce error: %s \n", err.Error())
		return false
	}

	return true
}
