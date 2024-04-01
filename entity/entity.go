package entity

import (
	"bytes"
	"encoding/gob"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
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
	Update()
	Register(entity IEntity)
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

func (e *Entity) Register(entity IEntity) {
	e.rType = reflect.TypeOf(entity)
	e.rVal = reflect.ValueOf(entity)
}

func (e *Entity) GetID() int64 {
	return e.ID
}

func (e *Entity) IsExistMethod(funcName string) bool {
	_, is := e.rType.MethodByName(funcName)
	return is
}

func (e *Entity) Call(packet pb.RpcPacket) {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	head := packet.Head
	v, ok := e.rType.MethodByName(head.FuncName)
	if !ok {
		logm.ErrorfE("方法不存在:%s\n", head.FuncName)
		return
	}

	// e.rVal.MethodByName("").Call()
	// 参数的个数
	nCount := v.Type.NumIn()

	//nCount1 := e.rVal.MethodByName(head.FuncName).Type().NumIn()

	//fmt.Println("-------------%d", nCount1)

	// 把所有的参数修改为valueof类型
	ps := make([]reflect.Value, nCount)

	buf := bytes.NewBuffer(packet.Buff)
	dec := gob.NewDecoder(buf)
	for i := 0; i < nCount; i++ {
		if i == 0 {
			ps[i] = e.rVal
			continue
		}
		// 获取每个参数的类型
		paramsValue := reflect.New(v.Type.In(i))

		dec.DecodeValue(paramsValue)

		ps[i] = paramsValue.Elem()
	}

	rets := v.Func.Call(ps)
	length := len(rets)
	for i := 0; i < length; i++ {
		logm.DebugfE("函数:%s 返回值: %v \n", rets[i])
	}
}

func (e *Entity) Update() {

}
