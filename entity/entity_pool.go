package entity

import (
	"github.com/mx5566/server/server/pb"
	"reflect"
	"sync"
)

type Opts struct {
}

type IEntityPool interface {
	AddEntity(entity IEntity)
	DelEntity(ID int64)
	GetEntity(ID int64) IEntity
	CallEntity(packet pb.RpcPacket)
}

type EntityPool struct {
	entity     IEntity // 这个entity存的池子里面具体是那种类型的实体
	entityMap  map[int64]IEntity
	entityLock *sync.RWMutex
}

func (p *EntityPool) InitPool(rType reflect.Type) {

	p.entityMap = make(map[int64]IEntity)

	// 存储池子里面的类型
	p.entity = reflect.New(rType).Interface().(IEntity)

	// 把类型注册到全局的类型管理器
	GEntityMgr.RegisterEntity(p.entity, WithType(EntityType_Pool), WithPool(p))
}

func (p *EntityPool) AddEntity(entity IEntity) {
	// 一个真实的实体
	entity.Register(entity)

	p.entityMap[entity.GetID()] = entity
}

func (p *EntityPool) DelEntity(ID int64) {
	delete(p.entityMap, ID)
}

func (p *EntityPool) GetEntity(ID int64) IEntity {
	if v, ok := p.entityMap[ID]; ok {
		return v
	}

	return nil
}

func (p *EntityPool) CallEntity(packet pb.RpcPacket) {
	ent := p.GetEntity(packet.Head.ID)
	if ent != nil {
		ent.Call(packet)
	}
}

func (p *EntityPool) Update() {
	for _, v := range p.entityMap {
		v.Update()
	}
}
