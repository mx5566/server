package gateserver

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/cluster"
	"github.com/mx5566/server/entity"
	"github.com/mx5566/server/network"
	"github.com/mx5566/server/rpc3"
	"github.com/mx5566/server/server/pb"
	"hash/crc32"
	"log"
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
	//s.Init()
	return s
}

func (p *ClientSession) SendToGameServer(funcName string, head rpc3.RpcHead, packet rpc3.RpcPacket) {
	head.DestServerType = rpc3.ServiceType_GameServer
	head.MsgSendType = rpc3.SendType_SendType_Single

	cluster.GCluster.SendMsg(head, packet)
}

func (p *ClientSession) SendToWorldServer(funcName string, head rpc3.RpcHead, packet rpc3.RpcPacket) {
	head.DestServerType = rpc3.ServiceType_WorldServer
	head.MsgSendType = rpc3.SendType_SendType_Single

	cluster.GCluster.SendMsg(head, packet)
}

func (p *ClientSession) Update() {
	/*socket := p.GetSocket()
	if socket == nil {
		return
	}

	if p.Session.DealQueue.IsEmpty() && !p.ReceQueue.IsEmpty() {
		p.DealQueue.Copy(p.ReceQueue)
	}

	if !p.DealQueue.IsEmpty() {
		for data := p.DealQueue.Pop(); data != nil; data = p.DealQueue.Pop() {
			p.HandlePacket(socket.GetConnId(), data)
		}
	}*/
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

	rpcPacket := pb.Marshal(head, &funcName, protoMsg)

	if head.DestServerType == rpc3.ServiceType_GameServer {
		//
		p.SendToGameServer(funcName, *head, rpcPacket)
	} else if head.DestServerType == rpc3.ServiceType_GateServer {
		// 加入是本地的话调用本地的方法
		// 我们需要根据类名 函数名 找到方法然后调用
		// 可以通过反射动态的获取方法，并且调用方法
		entity.GEntityMgr.Send(rpcPacket)
	} else if head.DestServerType == rpc3.ServiceType_LoginServer {

	}
}

func (p *ClientSession) HandleTest(ctx context.Context, test *pb.Test) {
	//TODO
	head := ctx.Value("rpcHead").(rpc3.RpcHead)
	// 转发到gameserver
	// 需要知道发送到那个服务器

	funcName := "AccountMgr.LoginAccountRequest"

	rpcPacket := pb.Marshal(&head, &funcName, test)
	//rpcPacketData, _ := proto.Marshal(&rpcPacket)
	//packet := rpc.Packet{
	//	Id:   head.ConnID,
	//	Buff: rpcPacketData, // RpcPacket 包含头和参数数据
	//}

	p.SendToGameServer(funcName, head, rpcPacket)

	log.Printf("接收测试数据 Name: %s, Password: %s\n", test.Name, test.PassWord)
}

func (p *ClientSession) HandleLoginAccount(ctx context.Context, msg *pb.LoginAccountReq) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	funcName := "AccountMgr.LoginAccountRequest"
	rpcPacket := pb.Marshal(&head, &funcName, msg)

	p.SendToWorldServer(funcName, head, rpcPacket)

}

func (p *ClientSession) HandleDisconnect(ctx context.Context, dis *pb.Disconnect) {
	log.Printf("客户端断开连接:%d\n", dis.ConnId)

}
