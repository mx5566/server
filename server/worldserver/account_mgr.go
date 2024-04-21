package worldserver

import (
	"context"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/pb"
)

var GAccountMgr = New()

type AccountMgr struct {
	entity.Entity
	cluster.ModuleAgent
}

func New() *AccountMgr {
	a := &AccountMgr{}
	return a
}

func (m *AccountMgr) Init() {
	m.Entity.Init()
	entity.RegisterEntity(m)
	m.Entity.Start()
	m.ModuleAgent.Init(rpc3.ModuleType_AccountMgr)
}

func (m *AccountMgr) RegisterAccount() {

}

func (m *AccountMgr) LoginAccountRequest(ctx context.Context, msg *pb.LoginAccountReq) {
	logm.DebugfE("账号登录请求:userName:%s, pass:%s", msg.GetUserName(), msg.GetPassword())
	// 返回一个消息
	packetHead := ctx.Value("rpcHead").(rpc3.RpcHead)

	// 处理数据库逻辑 记载账号获取账号相关数据
	cluster.GCluster.SendMsg(&rpc3.RpcHead{
		ClassName:      "",
		FuncName:       "",
		SrcServerID:    SERVER.GetID(),
		DestServerID:   packetHead.SrcServerID,
		DestServerType: rpc3.ServiceType_GateServer,
		ID:             0,
		ConnID:         packetHead.ConnID,
		MsgSendType:    rpc3.SendType_SendType_Single,
	}, "gateserver<-ClusterMsg.AccountLogin", 1, msg.GetUserName())

	cluster.GCluster.SendMsg(&rpc3.RpcHead{
		ClassName:      "",
		FuncName:       "",
		SrcServerID:    SERVER.GetID(),
		DestServerID:   packetHead.SrcServerID,
		DestServerType: rpc3.ServiceType_GateServer,
		ID:             0,
		ConnID:         packetHead.ConnID,
		MsgSendType:    rpc3.SendType_SendType_Single,
	}, "", "LoginAccontRep", &pb.LoginAccontRep{
		AccountId: 1,
	})

}

func (m *AccountMgr) OnModuleRegister() {
	logm.DebugfE("AccountMgr.OnModuleRegister----------------")
}
