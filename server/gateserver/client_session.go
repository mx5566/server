package gateserver

import (
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/entity"
	"github.com/mx5566/server/network"
	"github.com/mx5566/server/server/pb"
	"hash/crc32"
	"log"
	"reflect"
	"strings"
)

const (
	Send_Game = iota
	Send_Gate
	Send_Login
)

type ClientSession struct {
	network.Session
	gameID        uint32
	entity.Entity // 对象

}

func NewSession() network.ISession {
	s := &ClientSession{}
	s.Init()
	return s
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
	// 初始化实体
	p.Entity.Init()
	p.Session.Init()

	//p.RegisterPacket(1, network.HandleRegister{Handle: p.HandleAuth})
	//p.RegisterPacket(2, network.HandleRegister{Status: Send_Game})
	//p.RegisterPacket(2, network.HandleRegister{Status: Send_Game})

	p.RegisterPacketEx(&pb.Test{}, "gateserver<-ClientSession.HandleTest")

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
	_ = proto.Unmarshal(msg.MsgBody, protoMsg)

	// 接续函数
	// "gateserver<-ClientSession.HandleTest"
	funcName := route.FuncName
	strs := strings.Split(funcName, "<-")
	head := &pb.RpcHead{}
	if len(strs) == 2 {
		switch strs[0] {
		case "gateserver":
			head.DestServerType = Send_Gate
		case "gameserver":
			head.DestServerType = Send_Game
		case "loginserver":
			head.DestServerType = Send_Login
		}

		funcName = strs[1]
	}

	// 拆分类名和函数名
	strs = strings.Split(funcName, ".")
	if len(strs) == 2 {
		head.ClassName = strs[0] // 类名
		funcName = strs[1]       // 真正的函数名字
	}

	head.SrcServerID = SERVER.GetID()
	head.FuncName = funcName

	rpcPacket := pb.RpcPacket{}
	rpcPacket.Head = head

	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	enc.Encode(protoMsg)
	rpcPacket.Buff = buf.Bytes()

	if head.DestServerType == Send_Game {
		//
	} else if head.DestServerType == Send_Gate {
		// 加入是本地的话调用本地的方法
		// 我们需要根据类名 函数名 找到方法然后调用
		// 可以通过反射动态的获取方法，并且调用方法
		entity.GEntityMgr.Call(rpcPacket)
	} else if head.DestServerType == Send_Login {

	}

	/*
		if _, ok := p.Hrs[msg.MsgId]; ok {
			if p.Hrs[msg.MsgId].Handle != nil {
				p.Hrs[msg.MsgId].Handle(connId, msg)
			} else {
				// 没找到处理函数 根据状态发送到指定的服务器
				switch p.Hrs[msg.MsgId].Status {
				// gameserver游戏服务器
				case 1:
				case 2:
				default:
				}
			}
		}
	*/
}

func (p *ClientSession) HandleTest(test *pb.Test) {
	log.Printf("接收测试数据 Name: %s, Password: %s\n", test.Name, test.PassWord)
}

func (p *ClientSession) HandleAuth(connId uint32, msg *network.MsgPacket) bool {
	return true
}
