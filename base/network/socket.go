package network

import (
	"encoding/binary"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/rpc3"
	"log"
	"net"
)

const TcpHeadSize = 6
const TcpIDLength = 4

type HandleFunc func(packet rpc3.Packet) //回调函数

type ISocket interface {
	Init(ip string, port uint16) bool
	Start() bool
	Stop() bool
	Run() bool
	HandlePacket(rpc3.Packet)
	ReceivePacket([]byte) bool
	BindPacketFunc(HandleFunc)
	Connect() bool
	OnNetFail()
	Send(rpc3.Packet)
	GetConnId() uint32
	SetSessionType(SESSION_TYPE)
	GetSessionType() SESSION_TYPE
}

type Socket struct {
	Ip               string
	Port             uint16
	conn             net.Conn
	nReceBuff        []byte
	nMaxReceBuffSize int // 避免无效包
	nMaxSendBuffSize int
	connId           uint32
	session          ISession
	sessionType      SESSION_TYPE
	handleFunc       HandleFunc
}

func (s *Socket) GetConnId() uint32 {
	return s.connId
}

func (s *Socket) SetSessionType(sessionType SESSION_TYPE) {
	s.sessionType = sessionType
}

func (s *Socket) name() {

}

func (s *Socket) GetSessionType() SESSION_TYPE {
	return s.sessionType
}

func (s *Socket) Init(ip string, port uint16) bool {
	if s.Ip == ip && s.Port == port {
		return false
	}

	s.Ip = ip
	s.Port = port

	s.nMaxReceBuffSize = 1 * 1024 * 1024 // 10MB
	s.nMaxSendBuffSize = 10 * 1024 * 1024
	return true
}

func (s *Socket) Start() bool {

	return false
}

func (s *Socket) Stop() bool {
	if s.conn == nil {
		return true
	}

	s.conn.Close()

	return true
}

func (s *Socket) Send(rpcPacket rpc3.Packet) {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	if s.conn == nil {
		return
	} else if len(rpcPacket.Buff) > s.nMaxSendBuffSize {
		logm.PanicfE("send over maxsendbuff: %dMB\n", s.nMaxSendBuffSize/1024/1024)
	}

	_, _ = s.conn.Write(rpcPacket.Buff)
}

func (s *Socket) Run() bool {

	return true
}

func (s *Socket) BindPacketFunc(hFunc HandleFunc) {
	s.handleFunc = hFunc
}

func (s *Socket) HandlePacket(packet rpc3.Packet) {
	s.HandlePacket(packet)
}

func (s *Socket) ReceivePacket(data []byte) bool {
	// 因为可能有剩余的不够长度的数据包在里面，所以追加
	buff := append(s.nReceBuff, data...)
	s.nReceBuff = []byte{}
	curSize := 0

	// 包头占用四个字节 msgId + msgLen + msgBody
	// 粘包处理过程
	for {
		// 剩余的数据长度
		dataSize := len(buff[curSize:])

		if dataSize < TcpHeadSize {
			s.nReceBuff = buff[curSize:]
			return true
		}

		//
		msgLen := binary.BigEndian.Uint16(buff[curSize+TcpIDLength : curSize+TcpHeadSize])

		// 包的大小超过最大包的大小，无效包丢弃
		if int(TcpHeadSize+msgLen) > s.nMaxReceBuffSize {
			return false
		}

		// 消息的长度大于剩余的数据长度
		if int(TcpHeadSize+msgLen) > dataSize {
			s.nReceBuff = buff[curSize:]
			return true
		} else {
			packet := rpc3.Packet{
				Id:   s.connId,
				Buff: buff[curSize : uint16(curSize)+TcpHeadSize+msgLen],
			}
			s.handleFunc(packet)

			// 二进制解包
			//var dp DataPacket
			//msg := dp.Decode(buff[curSize : uint16(curSize)+TcpHeadSize+msgLen])
			// 放到消息队列里面
			//s.session.AddQueue(msg)
			curSize += int(TcpHeadSize + msgLen)
		}
	}

	return true
}

func (s *Socket) Connect() bool {
	var strRemote = fmt.Sprintf("%s:%d", s.Ip, s.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Printf("%v", err)
		return false
	}

	fmt.Println("------- ", strRemote)

	conn, err1 := net.DialTCP("tcp4", nil, tcpAddr)
	if err1 != nil {
		return false
	}

	conn.SetNoDelay(true)
	s.conn = conn

	fmt.Printf("连接成功，请输入信息！\n")
	return true
}

func (s *Socket) OnNetFail() {
	s.Stop()
}
