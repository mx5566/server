package network

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type Server struct {
	server  *ServerSocket
	playMgr *PlayerMgr
}

type Player struct {
	connId uint32
	ID     uint64
}

type PlayerMgr struct {
	players map[uint64]*Player
}

var SERVER Server

func handleLogin(connId uint32, msg *MsgPacket) bool {
	fmt.Println("调用 handleLogin 接收到客户端1001消息号")

	// 登录成功
	player := new(Player)
	player.connId = connId
	player.ID = 1

	SERVER.playMgr.players[player.ID] = player
	return true
}

func handleClientLost(connId uint32, msg *MsgPacket) bool {

	return true
}
func TestServerSocket(t *testing.T) {
	s := new(ServerSocket)
	s.Init("0.0.0.0", 8080)
	s.Start()

	//s.BindPacketFunc(1001, handleLogin)
	//s.BindPacketFunc(1, handleClientLost)

	SERVER.server = s
	SERVER.playMgr = new(PlayerMgr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	q := <-c

	fmt.Printf("server ------- signal:[%v]", q)
}
