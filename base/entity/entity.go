package entity

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/mpsc"
	"github.com/mx5566/server/base/rpc3"
	"reflect"
	"sync/atomic"
)

var gEntityID int64

type EntityType uint16

const (
	EntityType_Single EntityType = iota
	EntityType_Pool
)

type IEntity interface {
	Init()
	Start()
	Run()
	GetID() int64
	SetID(int64)
	IsExistMethod(funcName string) bool
	Call(rpc3.RpcPacket)
	Register(entity IEntity)
	GetEntityType() EntityType
	GetEntityPool() IEntityPool
	SetEntityPool(pool IEntityPool)
	SetEntityType(EntityType)
	Send(packet rpc3.RpcPacket)
	Self() IEntity
}

type Entity struct {
	ID         int64
	EntityName string
	rType      reflect.Type
	rVal       reflect.Value
	pool       IEntityPool
	entityType EntityType
	mailBox    *mpsc.Queue[*rpc3.RpcPacket]
	mailChan   chan bool
}

func (e *Entity) Init() {
	if e.ID == 0 {
		e.ID = atomic.AddInt64(&gEntityID, 1)
	}

	e.mailBox = mpsc.New[*rpc3.RpcPacket]()
	e.mailChan = make(chan bool)
}

func (e *Entity) Self() IEntity {
	return e
}

func (e *Entity) Start() {

	go e.Run()
}

func (e *Entity) Run() {
	for {
		switch {
		case <-e.mailChan:
			for data := e.mailBox.Pop(); data != nil; data = e.mailBox.Pop() {
				e.Call(*data)
			}
		}
	}
}

func (e *Entity) Send(packet rpc3.RpcPacket) {
	e.mailBox.Push(&packet)

	//logm.DebugfE("EntityCallFunName:%s", packet.Head.FuncName)

	e.mailChan <- true
}

func (e *Entity) SetID(iD int64) {
	e.ID = iD
}

func (e *Entity) SetEntityType(entityType EntityType) {
	e.entityType = entityType
}

func (e *Entity) GetEntityType() EntityType {
	return e.entityType
}

func (e *Entity) GetEntityPool() IEntityPool {
	return e.pool
}

func (e *Entity) SetEntityPool(pool IEntityPool) {
	e.pool = pool
}

func (e *Entity) Register(entity IEntity) {
	e.rType = reflect.TypeOf(entity)
	e.rVal = reflect.ValueOf(entity)
	e.EntityName = base.GetClassName(entity)
}

func (e *Entity) GetID() int64 {
	return e.ID
}

func (e *Entity) IsExistMethod(funcName string) bool {
	_, is := e.rType.MethodByName(funcName)
	return is
}

func (e *Entity) Call(packet rpc3.RpcPacket) {
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
	// 参数的个数
	nCount := v.Type.NumIn()

	// 把所有的参数修改为valueof类型
	ps := make([]reflect.Value, nCount)

	buf := bytes.NewBuffer(packet.Buff)
	dec := gob.NewDecoder(buf)
	for i := 0; i < nCount; i++ {
		if i == 0 {
			ps[i] = e.rVal
			continue
		}

		if i == 1 {
			ps[i] = reflect.ValueOf(context.WithValue(context.Background(), "rpcHead", *(*rpc3.RpcHead)(packet.Head)))
			continue
		}
		// 获取每个参数的类型
		paramsValue := reflect.New(v.Type.In(i))

		err := dec.DecodeValue(paramsValue)
		if err != nil {
			logm.ErrorfE("解析数据错误: %s", err.Error())
			continue
		}

		ps[i] = paramsValue.Elem()
	}

	//logm.DebugfE("Call:%s", packet.Head.FuncName)

	rets := v.Func.Call(ps)
	length := len(rets)
	for i := 0; i < length; i++ {
		logm.DebugfE("函数:%s 返回值: %v \n", rets[i])
	}
}
