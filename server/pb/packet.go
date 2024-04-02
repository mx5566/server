package pb

import (
	"bytes"
	"encoding/gob"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/network"
	"strings"
)

func Route(head *RpcHead, funcName string) string {
	strs := strings.Split(funcName, "<-")
	if len(strs) == 2 {
		switch strings.ToLower(strs[0]) {
		case "gateserver":
			head.DestServerType = network.Send_Gate
		case "gameserver":
			head.DestServerType = network.Send_Game
		case "loginserver":
			head.DestServerType = network.Send_Login
		}
		funcName = strs[1]
	}

	strs = strings.Split(funcName, ".")
	if len(strs) == 2 {
		head.ClassName = strs[0]
		head.FuncName = strs[1]
	}

	return funcName
}

func Marshal(head *RpcHead, funcName *string, params ...interface{}) RpcPacket {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	// 对函数结构分解
	// gameserver<-playermgr.Login
	*funcName = Route(head, *funcName)

	pac := RpcPacket{}

	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	length := len(params)
	for i := 0; i < length; i++ {
		enc.Encode(params[i])
	}
	pac.Head = head
	pac.Buff = buf.Bytes()

	return pac
}
