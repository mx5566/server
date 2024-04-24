package gateserver

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/pb"
	"hash/crc32"
	"reflect"
)

type Route struct {
	PmsgFunc func() proto.Message
	FuncName string
}

type ClientSession struct {
	gameID        uint32
	entity.Entity // 对象

	HrsStr map[string]Route
	HrsId  map[uint32]Route
}

func NewSession() *ClientSession {
	s := &ClientSession{}
	return s
}

func (p *ClientSession) SendToGameServer(head rpc3.RpcHead, funcName string, param ...interface{}) {
	head.DestServerType = rpc3.ServiceType_GameServer
	head.MsgSendType = rpc3.SendType_SendType_Single

	cluster.GCluster.SendMsg(&head, funcName, param...)
}

func (p *ClientSession) SendToWorldServer(head rpc3.RpcHead, funcName string, param ...interface{}) {
	head.DestServerType = rpc3.ServiceType_WorldServer
	head.MsgSendType = rpc3.SendType_SendType_Single

	cluster.GCluster.SendMsg(&head, funcName, param...)
}

func (p *ClientSession) Init() {

	p.HrsStr = make(map[string]Route)
	p.HrsId = make(map[uint32]Route)

	// 初始化实体

	p.RegisterPacket(&pb.Test{}, "gateserver<-ClientSession.HandleTest")
	p.RegisterPacket(&pb.Disconnect{}, "gateserver<-ClientSession.HandleDisconnect")
	p.RegisterPacket(&pb.LoginAccountReq{}, "gateserver<-ClientSession.HandleLoginAccount")

	p.Entity.Init()
	p.Entity.Start()
	entity.RegisterEntity(p)
}

func (p *ClientSession) RegisterPacket(msgName proto.Message, funcName string) {
	name := base.GetMessageName(msgName)

	packetFunc := func() proto.Message {
		val := reflect.ValueOf(msgName).Elem()
		e := reflect.New(val.Type())
		e.Elem().Field(3).Set(val.Field(3))

		return e.Interface().(proto.Message)
	}

	hr := Route{}
	hr.PmsgFunc = packetFunc
	hr.FuncName = funcName

	p.HrsStr[string(name)] = hr
	p.HrsId[crc32.ChecksumIEEE([]byte(name))] = hr
}

func (p *ClientSession) HandlePacket(packet rpc3.Packet) {
	connId := packet.Id
	buff := packet.Buff

	var dp network.DataPacket
	msg := dp.Decode(buff)

	// 根绝客户端的二进制消息,判断消息id是不是注册了
	if _, ok := p.HrsId[msg.MsgId]; !ok {
		logm.ErrorfE("错误的解析包 Id: %d \n", msg.MsgId)
		return
	}

	route := p.HrsId[msg.MsgId]

	// protobufmessage
	protoMsg := route.PmsgFunc() // 传递函数比传递一块内存节省空间

	head := &rpc3.RpcHead{}
	head.SrcServerID = SERVER.GetID()
	head.ConnID = connId
	head.ID = p.GetID()

	_ = proto.Unmarshal(msg.MsgBody, protoMsg)

	// 接续函数
	// "gateserver<-ClientSession.HandleTest"
	funcName := route.FuncName

	pb.Route(head, funcName)

	if head.DestServerType == rpc3.ServiceType_GameServer {
		//
		p.SendToGameServer(*head, funcName, protoMsg)
	} else if head.DestServerType == rpc3.ServiceType_GateServer {
		// 加入是本地的话调用本地的方法
		// 我们需要根据类名 函数名 找到方法然后调用
		// 可以通过反射动态的获取方法，并且调用方法
		//entity.GEntityMgr.Send(rpcPacket)

		entity.GEntityMgr.SendMsg(*head, funcName, protoMsg)
	} else if head.DestServerType == rpc3.ServiceType_LoginServer {
		p.SendToWorldServer(*head, funcName, protoMsg)
	}
}

func (p *ClientSession) HandleTest(ctx context.Context, test *pb.Test) {
	//TODO
	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	funcName := "AccountMgr.LoginAccountRequest"

	p.SendToGameServer(head, funcName, test)

	logm.DebugfE("接收测试数据 Name: %s, Password: %s\n", test.Name, test.PassWord)
}

func (p *ClientSession) HandleLoginAccount(ctx context.Context, msg *pb.LoginAccountReq) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	// 账号登录
	funcName := "AccountMgr.LoginAccountRequest"
	p.SendToWorldServer(head, funcName, msg)
}

func (p *ClientSession) HandleDisconnect(ctx context.Context, dis *pb.Disconnect) {
	logm.DebugfE("客户端断开连接:%d\n", dis.ConnId)

	//
	entity.GEntityMgr.SendMsg(rpc3.RpcHead{ConnID: dis.GetConnId()}, "PlayerMgr.AccountLogining", dis.GetConnId())

	// 通知到游戏服务器玩家下线了做业务处理
}
