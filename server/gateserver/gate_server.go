package gateserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/network"
	"time"
)

type GateServer struct {
}

var SERVER GateServer

func (gs *GateServer) Init() {
	// 日志初始化
	logm.Init("gateserver", map[string]string{"errFile": "gate_server.log", "logFile": "gate_server_error.log"}, "debug")
	s := new(network.ServerSocket)
	s.Init("0.0.0.0", 8080)
	s.Start()
	//s.BindPacketFunc()

}

func (gs *GateServer) Loop() {
	for {

		// 暂停20微淼
		time.Sleep(20 * time.Millisecond)
	}

}

// 可以用IP+PORT 求一个哈希值
func (gs *GateServer) GetID() uint32 {
	return 1
}
