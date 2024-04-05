package pb

import (
	"bytes"
	"encoding/gob"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/rpc3"
	"strings"
)

func Route(head *rpc3.RpcHead, funcName string) string {
	strs := strings.Split(funcName, "<-")
	if len(strs) == 2 {
		switch strings.ToLower(strs[0]) {
		case "gateserver":
			head.DestServerType = rpc3.ServiceType_GateServer
		case "gameserver":
			head.DestServerType = rpc3.ServiceType_GameServer
		case "loginserver":
			head.DestServerType = rpc3.ServiceType_LoginServer
		}
		funcName = strs[1]
	}

	strs = strings.Split(funcName, ".")
	if len(strs) == 2 {
		head.ClassName = strs[0]
		head.FuncName = strs[1]
	}

	return head.FuncName
}

func Marshal(head *rpc3.RpcHead, funcName *string, params ...interface{}) rpc3.RpcPacket {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	// 对函数结构分解
	// gameserver<-playermgr.Login
	*funcName = Route(head, *funcName)

	pac := rpc3.RpcPacket{}

	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	length := len(params)
	for i := 0; i < length; i++ {
		err := enc.Encode(params[i])
		if err != nil {
			logm.ErrorfE("gob 编码错误 那么：%s, err: %s", *funcName, err.Error())
		} else {
			//logm.DebugfE("编码成功")
		}
	}
	pac.Head = head
	pac.Buff = buf.Bytes()

	return pac
}
