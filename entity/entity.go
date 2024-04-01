package entity

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/server/pb"
	"reflect"
	"sync/atomic"
)

var gEntityID int64

type IEntity interface {
	Init()
	GetID() int64
	IsExistMethod(funcName string) bool
	Call(pb.RpcPacket)
}

type Entity struct {
	ID         int64
	EntityName string
	rType      reflect.Type
	rVal       reflect.Value
}

func (e *Entity) Init() {
	if e.ID == 0 {
		e.ID = atomic.AddInt64(&gEntityID, 1)
	}
}

func (e *Entity) GetID() int64 {
	return e.ID
}

func (e *Entity) IsExistMethod(funcName string) bool {
	_, is := e.rType.MethodByName(funcName)
	return is
}

func (e *Entity) Call(packet pb.RpcPacket) {
	head := packet.Head
	v, ok := e.rType.MethodByName(head.FuncName)
	if !ok {
		logm.ErrorfE("方法不存在:%s\n", head.FuncName)
		return
	}

	e.rVal.MethodByName("").Call()

	//v.Func.Call()
}
