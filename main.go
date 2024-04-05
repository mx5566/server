package main

import (
	"fmt"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/server/gameserver"
	"github.com/mx5566/server/server/gateserver"
	"github.com/mx5566/server/server/worldserver"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	server  *network.ServerSocket
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

func handleLogin(connId uint32, msg *network.MsgPacket) bool {
	fmt.Println("调用 handleLogin 接收到客户端1001消息号")

	// 登录成功
	player := new(Player)
	player.connId = connId
	player.ID = 1

	SERVER.playMgr.players[player.ID] = player
	return true
}

func handleClientLost(connId uint32, msg *network.MsgPacket) bool {

	return true
}
func main() {
	args := os.Args

	// 调用init
	if args[1] == "gate" {
		// init
		gateserver.SERVER.Init()

	} else if args[1] == "game" {
		gameserver.SERVER.Init()

	} else if args[1] == "world" {
		worldserver.SERVER.Init()

	}
	/*
		s := new(network.ServerSocket)
		s.Init("0.0.0.0", 8080)
		s.Start()

		s.BindPacketFunc(1001, handleLogin)
		s.BindPacketFunc(1, handleClientLost)

		SERVER.server = s
		SERVER.playMgr = new(PlayerMgr)
	*/

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	q := <-c

	// 退出处理

	fmt.Printf("server ------- signal:[%v]", q)
	fmt.Println("HelloWorld...")
}
