package network

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"io"
)

// 主动连接服务器的客户端

type IClientSocket interface {
	ISocket
}

type ClientSocket struct {
	Socket
}

func (s *ClientSocket) Init(ip string, port uint16) bool {
	ret := s.Socket.Init(ip, port)
	if !ret {
		return false
	}

	s.session = new(Session)
	s.session.Init()

	return true
}

func (s *ClientSocket) Start() bool {
	if s.Socket.Connect() {
		// 启动一个协成接收数据
		go s.Run()
	}

	return true
}

func (s *ClientSocket) Run() bool {
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
			logm.ErrorE("远程链接：%s已经关闭！\n", s.conn.RemoteAddr().String())
			return false
		}

		if err != nil {
			logm.ErrorE("数据接收错误：%s 错误: %s\n", s.conn.RemoteAddr().String(), err.Error())
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
