package network

import (
	"fmt"
	"github.com/mx5566/server/base"
	"io"
)

// 服务器接收的客户端连接的对象
type IServerSocketClient interface {
	ISocket
}

type ServerSocketClient struct {
	Socket
	sc *ServerSocket
}

func (s *ServerSocketClient) Init(ip string, port uint16) bool {
	s.Socket.Init(ip, port)

	factory := GSessionFactoryMgr.GetFactory(int(s.GetSessionType()))

	s.session = factory.CreateSession()
	s.session.SetSocket(s)
	s.session.SetFactory(factory)

	return true
}

func (s *ServerSocketClient) Start() bool {
	go s.Run()

	return true
}

func (s *ServerSocketClient) Run() bool {
	buff := make([]byte, 1024)

	// defer 的生命周期与函数绑定，函数返回他也跟随被释放
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				base.TraceCode(err)
			}
		}()

		// 连接不存在
		if s.conn == nil {
			return false
		}

		n, err := s.conn.Read(buff)
		if err == io.EOF {
			// 通知服务器自己断开了
			s.OnNetFail()

			fmt.Printf("远程链接：%s已经关闭！\n", s.conn.RemoteAddr().String())
			return false
		}

		if err != nil {
			s.OnNetFail()

			fmt.Printf("数据接收错误：%s 错误: %s\n", s.conn.RemoteAddr().String(), err.Error())
			return false
		}

		if n < 0 {
			return false
		}

		s.Socket.ReceivePacket(buff[0:n])

		return true
	}

	for {
		if !loop() {
			break
		}
	}

	return false
}

func (s *ServerSocketClient) OnNetFail() {
	// 连接断开了，需要通知到上层逻辑
	msg := new(MsgPacket)
	msg.MsgId = 1

	s.session.AddQueue(msg)

	s.Stop()

	s.sc.DelConn(s.connId)
}
