package account

import (
	"context"
	"github.com/mx5566/server/base/uuid"
	"github.com/mx5566/server/server/model"

	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/orm/mongodb"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/pb"
	"github.com/mx5566/server/server/worldserver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GAccountMgr = New()

type AccountMgr struct {
	entity.Entity
	cluster.ModuleAgent

	accounts map[int64]*Account

	mongoInstance *mongodb.MongoDB[AccountInfo]
	ctx           context.Context
}

func New() *AccountMgr {
	a := &AccountMgr{}
	return a
}

func (m *AccountMgr) Init() {
	m.mongoInstance = mongodb.NewMGDB[AccountInfo]("account", "account_tbl")
	m.ctx = context.Background()

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

	var errCode int32 = 0
	// 去accountdb 查找账号是否存在
	filter := mongodb.Newfilter().EQ("accountName", msg.UserName)

	ops := options.FindOneOptions{}
	ops.SetProjection(bson.D{{"accountId", 1}, {"accountPasswd", 1}})

	oneA, err := m.mongoInstance.FindOne(m.ctx, filter, &ops)
	if err != nil {
		logm.ErrorfE("账号数据库查询失败:%s", err.Error())
		errCode = 1
		//goto clientTag
	} else if oneA.AccountID == 0 {
		// 没找到，直接插入数据库
		ac := AccountInfo{
			AccountID:     uuid.UUID.UUID(),
			AccountName:   msg.GetUserName(),
			AccountPasswd: msg.GetPassword(),
		}
		result := m.mongoInstance.InsertOne(m.ctx, ac)
		if result == nil {
			errCode = 2
			//goto clientTag
		}
	} else if oneA.AccountID != 0 {
		if oneA.AccountPasswd != msg.GetPassword() {
			errCode = 3
			//goto clientTag
		}
	}

	if errCode != 0 {
		rep := pb.LoginAccontRep{AccountId: 0, ErrCode: errCode}
		cluster.GCluster.SendMsg(&rpc3.RpcHead{
			SrcServerID:    worldserver.SERVER.GetID(),
			DestServerID:   packetHead.SrcServerID,
			DestServerType: rpc3.ServiceType_GateServer,
			ID:             0,
			ConnID:         packetHead.ConnID,
		}, "", "LoginAccontRep", &rep)
		return
	}

	m.LoginAccount(packetHead, oneA.AccountID)
}

func (m *AccountMgr) loadAccount(accountID int64) *Account {
	// 去accountdb 查找账号是否存在
	filter := mongodb.Newfilter().EQ("accountId", accountID)

	ops := options.FindOneOptions{}
	ops.SetProjection(bson.D{{"_id", 0}})

	oneA, err := m.mongoInstance.FindOne(m.ctx, filter, &ops)
	if err != nil {
		return nil
	}

	account := new(Account)
	acInfo := AccountInfo{AccountID: accountID}
	acInfo.AccountName = oneA.AccountName
	account.accountInfo = acInfo
	account.roleInfos = m.loadPlayerSimple(accountID)

	return account
}

func (m *AccountMgr) loadPlayerSimple(accountID int64) (pls []*model.PlayerSimpleInfo) {
	filter := mongodb.Newfilter().EQ("accountID", accountID)
	pInstance := mongodb.NewMGDB[model.PlayerSimpleInfo]("game", "player_tbl")

	ops := options.FindOneOptions{}
	ops.SetProjection(bson.D{{"_id", 0}})

	// 查找角色列表 查找角色
	players, err := pInstance.Find(m.ctx, filter, 0)
	if err != nil {
		logm.ErrorfE("数据库查找角色失败:%s", err.Error())
		return
	}

	for _, v := range players {
		player := v
		pls = append(pls, &player)
	}

	return
}

func (m *AccountMgr) LoginAccount(packetHead rpc3.RpcHead, accountID int64) {
	account := m.GetAccount(accountID)
	if account == nil {
		account = m.loadAccount(accountID)
		m.accounts[accountID] = account
	}

	account.GateClusterID = packetHead.SrcServerID

	rep := pb.LoginAccontRep{AccountId: accountID}
	rep.ErrCode = 0
	for _, v := range account.roleInfos {
		pl := pb.PlayerList{}
		pl.Level = v.Level
		pl.Gold = v.Gold
		pl.PlayerId = v.PlayerID
		pl.PlayerName = v.Name
		pl.AccountID = v.AccountID
		rep.PList = append(rep.PList, &pl)
	}

	cluster.GCluster.SendMsg(&rpc3.RpcHead{
		SrcServerID:    worldserver.SERVER.GetID(),
		DestServerID:   packetHead.SrcServerID,
		DestServerType: rpc3.ServiceType_GateServer,
		ID:             0,
		ConnID:         packetHead.ConnID,
	}, "", "LoginAccontRep", &rep)

}

func (m *AccountMgr) LoginPlayerRequest(ctx context.Context, msg *pb.LoginPlayerReq) {
	logm.DebugfE("角色登录请求: %d", msg.PlayerId)
	packetHead := ctx.Value("rpcHead").(rpc3.RpcHead)

	accountID := msg.GetAccountID()
	playerID := msg.GetPlayerId()

	account := m.GetAccount(accountID)
	if account != nil {
		ret := account.PlayerLogin(playerID)
		if !ret {
			logm.ErrorfE("登录的角色不存: %d, %d", playerID, accountID)
			// 发个消息给客户端
			return
		}

		// 去gameserver
		cluster.GCluster.SendMsg(&rpc3.RpcHead{
			SrcServerID:    worldserver.SERVER.GetID(),
			DestServerID:   0,
			DestServerType: rpc3.ServiceType_GameServer,
			ConnID:         packetHead.GetConnID(),
		}, "gameserver<-PlayerMgr.PlayerLoginRequest", accountID, playerID, account.GateClusterID)
	}
}

func (m *AccountMgr) GetAccount(aID int64) *Account {
	if _, ok := m.accounts[aID]; ok {
		return m.accounts[aID]
	}
	return nil
}

func (m *AccountMgr) OnModuleRegister() {
	logm.DebugfE("AccountMgr.OnModuleRegister----------------")
}
