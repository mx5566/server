package gateserver

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/entity"
	"github.com/mx5566/server/network"
	"github.com/mx5566/server/server/pb"
	"hash/crc32"
	"log"
	"reflect"
)

type ClientSession struct {
	network.Session
	gameID        uint32
	entity.Entity // 对象

}

func NewSession() network.ISession {
	s := &ClientSession{}
	//s.Init()
	return s
}

func (p *ClientSession) SendToGameServer(funcName string, head pb.RpcHead, packet pb.Packet) {
	head.DestServerType = network.Send_Game

}

func (p *ClientSession) Update() {
	socket := p.GetSocket()
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
	}
}

func (p *ClientSession) Init() {
	p.SetID(int64(p.GetSocket().GetConnId()))
	// 初始化实体
	p.Entity.Init()
	p.Session.Init()

	//p.RegisterPacket(1, network.HandleRegister{Handle: p.HandleAuth})
	//p.RegisterPacket(2, network.HandleRegister{Status: Send_Game})
	//p.RegisterPacket(2, network.HandleRegister{Status: Send_Game})

	p.RegisterPacketEx(&pb.Test{}, "gateserver<-ClientSession.HandleTest")
	p.RegisterPacketEx(&pb.Disconnect{}, "gateserver<-ClientSession.HandleDisconnect")

	GSessionMgr.AddSession(p)
}

func (p *ClientSession) RegisterPacketEx(msgName proto.Message, funcName string) {
	name := base.GetMessageName(msgName)

	packetFunc := func() proto.Message {
		val := reflect.ValueOf(msgName).Elem()
		e := reflect.New(val.Type())
		e.Elem().Field(3).Set(val.Field(3))

		return e.Interface().(proto.Message)
	}

	hr := network.Route{}
	hr.PmsgFunc = packetFunc
	hr.FuncName = funcName

	p.Session.HrsStr[string(name)] = hr
	p.Session.HrsId[crc32.ChecksumIEEE([]byte(name))] = hr
}

func (p *ClientSession) RegisterPacket(msgId uint32, hr network.HandleRegister) {
	p.Hrs[msgId] = hr
}

func (p *ClientSession) HandlePacket(connId uint32, msg *network.MsgPacket) {
	if msg == nil {
		return
	}

	// 根绝客户端的二进制消息,判断消息id是不是注册了
	if _, ok := p.Session.HrsId[msg.MsgId]; !ok {
		logm.ErrorfE("错误的解析包 Id: %d \n", msg.MsgId)
		return
	}

	route := p.Session.HrsId[msg.MsgId]

	// protobufmessage
	protoMsg := route.PmsgFunc() // 传递函数比传递一块内存节省空间

	head := &pb.RpcHead{}
	head.SrcServerID = SERVER.GetID()
	head.ConnID = connId
	head.ID = p.GetID()

	_ = proto.Unmarshal(msg.MsgBody, protoMsg)

	// 接续函数
	// "gateserver<-ClientSession.HandleTest"
	funcName := route.FuncName

	rpcPacket := pb.Marshal(head, &funcName, protoMsg)

	if head.DestServerType == network.Send_Game {
		//

	} else if head.DestServerType == network.Send_Gate {
		// 加入是本地的话调用本地的方法
		// 我们需要根据类名 函数名 找到方法然后调用
		// 可以通过反射动态的获取方法，并且调用方法
		entity.GEntityMgr.Call(rpcPacket)
	} else if head.DestServerType == network.Send_Login {

	}
}

func (p *ClientSession) HandleTest(ctx context.Context, test *pb.Test) {
	//TODO
	head := ctx.Value("rpcHead").(pb.RpcHead)
	// 转发到gameserver
	// 需要知道发送到那个服务器

	funcName := "AccountMgr.LoginAccountRequest"

	rpcPacket := pb.Marshal(&head, &funcName, test)
	rpcPacketData, _ := proto.Marshal(&rpcPacket)
	packet := pb.Packet{
		Id:   head.ConnID,
		Buff: rpcPacketData, // RpcPacket 包含头和参数数据
	}

	p.SendToGameServer(funcName, head, packet)

	log.Printf("接收测试数据 Name: %s, Password: %s\n", test.Name, test.PassWord)
}

func (p *ClientSession) HandleAuth(connId uint32, msg *network.MsgPacket) bool {

	return true
}

func (p *ClientSession) HandleDisconnect(ctx context.Context, dis *pb.Disconnect) {
	log.Printf("客户端断开连接:%d\n", dis.ConnId)

	GSessionMgr.DelSession(p.GetID())

}
