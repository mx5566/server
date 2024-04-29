package player

import (
	"context"
	"fmt"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/aoi"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"sync"
)

type IUnit interface {
	OnEnterAoi(aoi *aoi.AOI)
	OnLeaveAoi(aoi *aoi.AOI)
}

type Unit struct {
	base.Vector3
	I        IUnit // 实体对象
	Aoi      aoi.AOI
	ID       int64
	TypeName string
	Attr     map[string]int64 // 属性
	MapID    uint32

	InterestedIn sync.Map
	InterestedBy sync.Map
}

func (u *Unit) Init() {

}

func (u *Unit) SendUnitStatus(ctx context.Context, playerID int64, gateClusterID uint32) {
	fmt.Println("发送状态信息给", playerID)
}

// aoi 对象的可是数据发送给指定的玩家
func (u *Unit) SendRemoteData(ctx context.Context, playerID int64, gateClusterID uint32) {
	fmt.Println("发送远程信息给", playerID)

}

func (u *Unit) AddInterestIn(ctx context.Context, unit *Unit) {
	u.InterestedIn.Store(unit, struct {
	}{})

	// 对方状态、属性通知给自己
	// 判断自己是不是角色 是角色就发消息给自己

	// 获取 网关的ID 玩家ID
	if u.TypeName == "Player" {
		// player := u.I.(*Player)
		// 第一个消息通知对方的位置状态给自己
		entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: unit.ID}, unit.TypeName+".SendUnitStatus", u.ID /*接收数据玩家的ID*/, 1001 /*接收数据的玩家网关ID*/)

		// 第二个消息把对方的可以被看到的数据发送给自己
		entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: unit.ID}, unit.TypeName+".SendRemoteData", u.ID /*接收数据玩家的ID*/, 1001 /*接收数据的玩家网关ID*/)
	}
}

func (u *Unit) AddInterestBy(ctx context.Context, unit *Unit) {
	u.InterestedBy.Store(unit, struct {
	}{})

	// 对方状态、属性通知给自己
	// 判断自己是不是角色 是角色就发消息给自己
}

func (u *Unit) UnInterestIn(ctx context.Context, unit *Unit) {
	u.InterestedIn.Delete(unit)

	// 从视野列表删除对方
	// 判断自己是不是角色 是角色就发消息给自己
}

func (u *Unit) UnInterestBy(ctx context.Context, unit *Unit) {
	u.InterestedBy.Delete(unit)

	// 从被观察列表删除对方
	// 判断自己是不是角色 是角色就发消息给自己

}

func (u *Unit) OnEnterAoi(aoi *aoi.AOI) {
	u1 := aoi.EntityData.(*Unit)

	u.InterestedIn.Store(u1, struct{}{})

	u1.InterestedBy.Store(u, struct{}{})

	//entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: u.ID}, u.TypeName+".AddInterestIn", u1)

	//entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: u1.ID}, u.TypeName+".AddInterestBy", u)

	fmt.Println("OnEnterAoi--------------------")
}

func (u *Unit) OnLeaveAoi(aoi *aoi.AOI) {
	u1 := aoi.EntityData.(*Unit)

	entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: u.ID}, u.TypeName+".UnInterestIn", u1)

	entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: u1.ID}, u.TypeName+".UnInterestBy", u)
}
