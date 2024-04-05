package entity

import (
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/rpc3"
	"github.com/mx5566/server/server/pb"
)

var GEntityMgr = CreateEntityMgr()

type Op struct {
	name       string //name
	entityType EntityType
	pool       IEntityPool
}

type OpOption func(*Op)

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func (op *Op) IsActorType(entityType EntityType) bool {
	return op.entityType == entityType
}

func WithType(entityType EntityType) OpOption {
	return func(op *Op) {
		op.entityType = entityType
	}
}

func WithPool(pPool IEntityPool) OpOption {
	return func(op *Op) {
		op.pool = pPool
	}
}

type EntityMgr struct {
	Entitys map[string]IEntity // 全局注册 一个类只有个注册，不管有多少个对象
}

func CreateEntityMgr() *EntityMgr {
	return &EntityMgr{Entitys: map[string]IEntity{}}
}

func (m *EntityMgr) SendMsg(head *rpc3.RpcHead, funcName string, params ...interface{}) {
	head.ConnID = 0
	rpcPacket := pb.Marshal(head, &funcName, params...)

	m.Send(rpcPacket)
}

// 注册是全局的 只有一次注册
// 对于需要多次注册的实体 比如玩家作为实体注册，我们不会把全部的玩家都注册进来因为我们调用其实只是获取实体类的方法，而不是单个对象的方法
// 但是我们有需要保存所有的实体对象
func (m *EntityMgr) RegisterEntity(entity IEntity, params ...OpOption) {
	op := Op{}
	op.applyOpts(params)

	name := base.GetClassName(entity)

	if _, ok := m.Entitys[name]; ok {
		logm.PanicfE("重复的注册实体Name: %s \n", name)
		return
	}

	// 记录实体名字和实体的映射
	m.Entitys[name] = entity

	entity.Register(entity)
	if op.pool != nil {
		entity.SetEntityPool(op.pool)
	}

	entity.SetEntityType(op.entityType)
}

func (m *EntityMgr) Send(packet rpc3.RpcPacket) {
	className := packet.Head.ClassName
	funcName := packet.Head.FuncName
	if v, ok := m.Entitys[className]; ok && v != nil {
		if v.IsExistMethod(funcName) {
			switch v.GetEntityType() {
			case EntityType_Single:
				v.Send(packet)
			case EntityType_Pool:
				v.GetEntityPool().CallEntity(packet)
			}
		}
	}
}

func (m *EntityMgr) Call(packet rpc3.RpcPacket) {
	className := packet.Head.ClassName
	funcName := packet.Head.FuncName
	if v, ok := m.Entitys[className]; ok && v != nil {
		if v.IsExistMethod(funcName) {
			switch v.GetEntityType() {
			case EntityType_Single:
				v.Call(packet)
			case EntityType_Pool:
				v.GetEntityPool().CallEntity(packet)
			}
		}
	}
}

func (m *EntityMgr) PacketFunc(packet rpc3.Packet) {
	rpcPacket := &rpc3.RpcPacket{}
	proto.Unmarshal(packet.Buff, rpcPacket)
	m.Send(*rpcPacket)
}

func RegisterEntity(entity IEntity) {
	GEntityMgr.RegisterEntity(entity)
}
