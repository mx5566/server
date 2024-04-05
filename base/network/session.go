package network

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/server/base"
)

type HandleRegister struct {
	Status uint16 // 转发到哪个服务器
	Handle HandleFunc
}

type Route struct {
	PmsgFunc func() proto.Message
	FuncName string
}

type ISession interface {
	Init()
	Update()
	AddQueue(msg *MsgPacket)
	SetSocket(socket ISocket)
	GetSocket() ISocket
	SetFactory(factory ISessionFactory)
}

type Session struct {
	ReceQueue     *base.FastQueue[*MsgPacket]
	DealQueue     *base.FastQueue[*MsgPacket]
	socket        ISocket
	packetHandles map[uint32]HandleFunc
	Hrs           map[uint32]HandleRegister
	factory       ISessionFactory
	HrsStr        map[string]Route
	HrsId         map[uint32]Route
}

func (s *Session) SetFactory(factory ISessionFactory) {
	s.factory = factory
}

func (s *Session) SetSocket(socket ISocket) {
	s.socket = socket
}

func (s *Session) GetSocket() ISocket {
	return s.socket
}

func (s *Session) Init() {
	s.ReceQueue = base.CreateFastQueue[*MsgPacket](true)
	s.DealQueue = base.CreateFastQueue[*MsgPacket](false)
	s.packetHandles = make(map[uint32]HandleFunc)

	s.Hrs = make(map[uint32]HandleRegister)
	s.HrsStr = make(map[string]Route)
	s.HrsId = make(map[uint32]Route)
}

func (s *Session) AddQueue(msg *MsgPacket) {
	s.ReceQueue.Push(msg)
}

// BindPacketFunc 绑定消息与处理函数
func (s *Session) BindPacketFunc(msgId uint32, handle HandleFunc) {
	if _, ok := s.packetHandles[msgId]; ok {
		fmt.Printf("重复的绑定一个消息: %d", msgId)
		return
	}

	s.packetHandles[msgId] = handle
}

// HandlePacket 处理消息
func (s *Session) HandlePacket(connId uint32, msg *MsgPacket) {
	if msg == nil {
		return
	}

	if _, ok := s.packetHandles[msg.MsgId]; ok {
		//s.packetHandles[msg.MsgId](connId, msg)
	}

}

// Update 每帧更新
func (s *Session) Update() {
	if s.socket == nil {
		return
	}

	if !s.DealQueue.IsEmpty() && s.ReceQueue.IsEmpty() {
		s.DealQueue.Copy(s.ReceQueue)
	}

	if !s.DealQueue.IsEmpty() {
		for data := s.DealQueue.Pop(); data != nil; data = s.DealQueue.Pop() {
			s.HandlePacket(s.socket.GetConnId(), data)
		}
	}
}
