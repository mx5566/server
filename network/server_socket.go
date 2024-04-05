package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// 服务器listening
type IServerSocket interface {
	ISocket
}

type ServerSocket struct {
	Socket
	ln          *net.TCPListener
	Clients     map[uint32]*ServerSocketClient
	clientMutex sync.Mutex
	rndId       uint32
}

func (s *ServerSocket) Init(ip string, port uint16) bool {
	s.Clients = make(map[uint32]*ServerSocketClient)
	s.Socket.Init(ip, port)
	return true
}

func (s *ServerSocket) Start() bool {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addt error : ", err)
		return false
	}

	//2 监听服务器的地址
	ln, err1 := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println("listen ", "tcp4", " err ", err1)
		return false
	}

	s.ln = ln

	go s.Run()

	fmt.Println("Start server ok...")

	return true
}

func (s *ServerSocket) DelConn(ConnId uint32) {
	s.clientMutex.Lock()
	delete(s.Clients, ConnId)
	s.clientMutex.Unlock()
}

func (s *ServerSocket) AddConn(conn *net.TCPConn) {
	ssc := new(ServerSocketClient)
	//ssc.SetSessionType(s.GetSessionType())
	ssc.BindPacketFunc(s.handleFunc)

	barray := strings.Split(conn.RemoteAddr().String(), ":")
	ret, _ := strconv.Atoi(barray[1])
	ssc.Init(barray[0], uint16(ret))

	ssc.conn = conn
	ssc.connId = atomic.AddUint32(&s.rndId, 1)
	ssc.sc = s

	s.clientMutex.Lock()
	s.Clients[ssc.connId] = ssc
	s.clientMutex.Unlock()

	ssc.Start()
}

func (s *ServerSocket) Stop() bool {
	//
	s.clientMutex.Lock()
	for key, _ := range s.Clients {
		client := s.Clients[key]
		client.Stop()

	}
	s.clientMutex.Unlock()
	s.ln.Close()

	return false
}

func (s *ServerSocket) Run() bool {
	for {
		conn, err := s.ln.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err", err)
			continue
		}

		s.AddConn(conn)

	}

	return false
}
