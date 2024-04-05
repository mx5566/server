package worldserver

import (
	"context"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/server/pb"
)

var GAccountMgr = New()

type AccountMgr struct {
	entity.Entity
}

func New() *AccountMgr {
	a := &AccountMgr{}

	a.Init()

	return a
}

func (m *AccountMgr) Init() {
	m.Entity.Init()
	m.Entity.Start()
	entity.RegisterEntity(m)
}

func (m *AccountMgr) RegisterAccount() {

}

func (m *AccountMgr) LoginAccountRequest(ctx context.Context, msg *pb.LoginAccountReq) {
	logm.DebugfE("账号登录请求:userName:%s, pass:%s", msg.GetUserName(), msg.GetPassword())
}
