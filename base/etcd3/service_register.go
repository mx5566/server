package etcd3

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/rpc3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// 服务注册模块
type ServiceRegister struct {
	client    *clientv3.Client
	lease     clientv3.Lease
	info      *rpc3.ClusterInfo
	leaseID   clientv3.LeaseID
	timeGrant int64
	status    int // 0  1
}

func NewServiceRegister(clusterInfo *rpc3.ClusterInfo, config conf.ServiceEtcd) *ServiceRegister {
	s := &ServiceRegister{}
	s.Init(clusterInfo, config.EndPoints, config.GrantTime)

	return s
}

func (r *ServiceRegister) Init(clusterInfo *rpc3.ClusterInfo, endPoints []string, timeNum int64) {
	conf := clientv3.Config{
		Endpoints: endPoints,
		//DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(conf)
	if err != nil {
		logm.PanicfE("服务注册模块启动失败: %v, 错误: %s\n", endPoints, err.Error())
		return
	}

	lease := clientv3.NewLease(client)

	r.client = client
	r.lease = lease
	r.timeGrant = timeNum
	r.status = 0
	r.info = clusterInfo

	go r.Run()
}

// 设置续租时间
func (r *ServiceRegister) SetLease() {
	//设置租约时间
	leaseResp, err := r.lease.Grant(context.Background(), r.timeGrant)
	if err != nil {
		logm.ErrorfE("设置租约时间错误: %s\n", err.Error())
		return
	}

	r.leaseID = leaseResp.ID

	key := base.ServiceName + r.info.ServiceType.String() + "/" + fmt.Sprintf("%s:%d", r.info.Ip, r.info.Port)
	val, _ := json.Marshal(r.info)
	_, err = r.client.Put(context.Background(), key, string(val), clientv3.WithLease(r.leaseID))
	if err != nil {
		logm.ErrorfE("etcd lease Put error: %s \n", err.Error())
		return
	}

	// 进入续约状态
	r.status = 1
}

func (r *ServiceRegister) KeepAlive() {
	_, err := r.lease.KeepAliveOnce(context.Background(), r.leaseID)
	if err != nil {
		r.status = 0
		logm.ErrorfE("etcd lease KeepAliveOnce error: %s \n", err.Error())
		return
	}

	// 避免cpu忙
	time.Sleep(time.Duration(r.timeGrant / 2))
}

// 监听 续租情况
func (r *ServiceRegister) Run() {
	for {
		switch r.status {
		// 初始状态
		case 0:
			r.SetLease()
		// 续约
		case 1:
			r.KeepAlive()
		}
	}
}
