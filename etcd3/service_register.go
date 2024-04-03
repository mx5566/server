package etcd3

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/rpc"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 服务注册模块
type ServiceRegister struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	info          *rpc.ClusterInfo
	leaseID       clientv3.LeaseID
	timeGrant     int64
	status        int // 0  1
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	cancelFunc    func()
}

func NewServiceRegister(endPoints []string, timeNum int64) *ServiceRegister {
	s := &ServiceRegister{}
	s.Init(endPoints, timeNum)

	return s
}

func (r *ServiceRegister) Init(endPoints []string, timeNum int64) {
	conf := clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(conf)
	if err != nil {
		logm.PanicfE("服务注册模块启动失败: %v, 错误: %s\n", endPoints, err.Error())
		return
	}

	lease := clientv3.NewLease(r.client)

	r.client = client
	r.lease = lease
	r.timeGrant = timeNum
	r.status = 0

	r.SetLease()

	go r.Run()
}

// 设置续租时间
func (r *ServiceRegister) SetLease() {
	//设置租约时间
	leaseResp, err := r.lease.Grant(context.Background(), r.timeGrant)
	if err != nil {
		return
	}

	r.leaseID = leaseResp.ID

	key := ServiceName + string(r.info.ServiceType) + "/" + fmt.Sprintf("%s:%d", r.info.Ip, r.info.Port)
	val, _ := json.Marshal(r.info)
	_, err = r.client.Put(context.Background(), key, string(val), clientv3.WithLease(r.leaseID))
	if err != nil {
		logm.PanicfE("etcd lease Put error: %s \n", err.Error())
		return
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	r.keepAliveChan, err = r.lease.KeepAlive(ctx, r.leaseID)
	if err != nil {
		logm.ErrorfE("etcd lease KeepAlive error: %s \n", err.Error())
		return
	}

	r.cancelFunc = cancelFunc
}

// 监听 续租情况
func (r *ServiceRegister) Run() {
	for {
		select {
		case leaseKeepResp := <-r.keepAliveChan:
			if leaseKeepResp == nil {
				logm.InfofE("已经关闭续租功能\n")
				return
			} else {
				logm.InfofE("续租成功\n")
			}
		}
	}
}

func (r *ServiceRegister) RevokeLease() {
	var err error
	go func(err error) {
		r.cancelFunc()
		time.Sleep(2 * time.Second)
		_, err = r.lease.Revoke(context.TODO(), r.leaseID)
	}(err)
}
