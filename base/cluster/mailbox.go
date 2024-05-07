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

// 信箱
type MailBox struct {
	mailBoxs map[int64]*rpc3.MailBox
	sync.Mutex
	client    *clientv3.Client
	lease     clientv3.Lease
	timeGrant int64
}

func (b *MailBox) Init(endPoints []string, t int64) {
	conf := clientv3.Config{
		Endpoints:         endPoints,
		DialKeepAliveTime: 10,
		//DialTimeout: time.Duration(t) * time.Second,
	}

	client, err := clientv3.New(conf)
	if err != nil {
		logm.PanicfE("MailBox连接etcd3模块启动失败: %v, 错误: %s\n", endPoints, err.Error())
		return
	}
	lease := clientv3.NewLease(client)

	b.client = client
	b.lease = lease
	b.timeGrant = t

	// 初始化
	b.mailBoxs = make(map[int64]*rpc3.MailBox)

	go b.Run()
}

func (b *MailBox) Register(mailBox *rpc3.MailBox, agent *MailBoxAgent) bool {
	//设置租约时间
	leaseResp, err := b.lease.Grant(context.Background(), b.timeGrant)
	if err != nil {
		logm.ErrorfE("设置租约时间错误: %s\n", err.Error())
		return false
	}

	agent.leaseID = leaseResp.ID

	key := base.MailBoxDir + fmt.Sprintf("%d", mailBox.GetID())
	val, _ := json.Marshal(mailBox)

	tx := b.client.Txn(context.Background())

	// CAS
	logm.DebugfE("MailBox提交事务的数据:%s, %s", key, string(val))
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(clientv3.OpPut(key, string(val), clientv3.WithLease(leaseResp.ID))).Else()

	resp, err := tx.Commit()
	if err != nil {
		logm.ErrorfE("MailBoxetcd3提交模块事务失败: %s %s",
			err.Error(),
			mailBox.String())
		return false
	}

	if !resp.Succeeded {
		logm.ErrorfE("MailBoxetcd3提交模块事务处理失败: %d %s %d",
			mailBox.GetID(),
			mailBox.GetMType().String(),
			mailBox.GetClusterID())
		return false
	}

	logm.DebugfE("事务提交数据: %d %s %d",
		mailBox.GetID(),
		mailBox.GetMType().String(),
		mailBox.GetClusterID())

	return true
}

func (b *MailBox) Run() {
	watchChan := b.client.Watch(context.Background(), base.MailBoxDir, clientv3.WithPrefix(), clientv3.WithPrevKV())
	b.GetAll()
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type.String() {
			case "PUT":
				b.AddMailBox(event.Kv.Value)
			case "DELETE":
				//logm.DebugfE("ModuleMgr Run DeleModule: %s %s cv:%d, mv:%d, vv:%d, ls: %d",
				//	string(event.Kv.Key), string(event.Kv.Value),
				//	event.Kv.CreateRevision, event.Kv.ModRevision,
				//	event.Kv.Version, event.Kv.Lease)
				b.DelMailBox(event.PrevKv.Value)
			}
		}
	}

	return
}

func (b *MailBox) GetAll() {
	resp, err := b.client.Get(context.Background(), base.MailBoxDir, clientv3.WithPrefix())
	if err != nil {
		return
	}
	i := 0
	for _, kv := range resp.Kvs {
		i++
		b.AddMailBox(kv.Value)
	}
}

func (b *MailBox) AddMailBox(data []byte) {
	mailBox := &rpc3.MailBox{}
	err := json.Unmarshal(data, mailBox)
	if err != nil {
		logm.ErrorfE("增加信箱数据解析失败:%s", err.Error())
		return
	}

	b.Lock()
	defer b.Unlock()
	b.mailBoxs[mailBox.GetID()] = mailBox

	logm.DebugfE("添加信箱成功:%d %s %d", mailBox.GetID(), mailBox.GetMType().String(), mailBox.GetClusterID())
}

func (b *MailBox) DelMailBox(data []byte) {
	mailBox := &rpc3.MailBox{}
	err := json.Unmarshal(data, mailBox)
	if err != nil {
		logm.ErrorfE("删除信箱数据解析失败:%s", err.Error())
		return
	}

	logm.DebugfE("删除信箱成功:%d %s %d", mailBox.GetID(), mailBox.GetMType().String(), mailBox.GetClusterID())

	b.Lock()
	defer b.Unlock()
	delete(b.mailBoxs, mailBox.GetID())
}

func (b *MailBox) GetMailBox(playerId int64) *rpc3.MailBox {
	b.Mutex.Lock()
	v, ok := b.mailBoxs[playerId]
	b.Mutex.Unlock()

	if ok {
		return v
	}
	return nil
}

func (b *MailBox) Lease(agent *MailBoxAgent) bool {
	_, err := b.lease.KeepAliveOnce(context.Background(), agent.leaseID)
	if err != nil {
		logm.ErrorfE("etcd mailbox lease KeepAliveOnce error: %s", err.Error())
		return false
	}
	return true
}

func (b *MailBox) Delete(id int64) bool {
	key := base.MailBoxDir + fmt.Sprintf("%d", id)
	_, err := b.client.Delete(context.Background(), key)
	if err != nil {
		logm.ErrorfE("etcd mailbox delete KeepAliveOnce error: %s", err.Error())
		return false
	}
	return true
}
