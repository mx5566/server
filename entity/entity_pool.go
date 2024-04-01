package entity

import (
	"reflect"
	"sync"
)

type IEntityPool interface {
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
	GEntityMgr.RegisterEntity(p.entity)
}

func (p *EntityPool) AddEntity(entity IEntity) {
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

func (p *EntityPool) Update() {
	for _, v := range p.entityMap {
		v.Update()
	}
}
