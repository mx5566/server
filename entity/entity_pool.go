package entity

import "sync"

type IEntityPool interface {
}

type EntityPool struct {
	entity     IEntity // 这个entity存的池子里面具体是那种类型的实体
	entityMap  map[int64]IEntity
	entityLock *sync.RWMutex
}

func (p *EntityPool) InitPool() {

	p.entityMap = make(map[int64]IEntity)
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
