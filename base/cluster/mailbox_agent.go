package cluster

import (
	"github.com/mx5566/server/base/rpc3"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type MailBoxAgent struct {
	mailBox rpc3.MailBox
	leaseID clientv3.LeaseID
}

func (m *MailBoxAgent) Init(t rpc3.MailBox) {
	m.mailBox = t
}

func (m *MailBoxAgent) RegisterAgent() bool {
	return GCluster.mailBox.Register(&m.mailBox, m)
}

func (m *MailBoxAgent) Lease() {
	GCluster.mailBox.Lease(m)
}

func (m *MailBoxAgent) Delete() {
	GCluster.mailBox.Delete(m.mailBox.GetID())
}
