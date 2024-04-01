package gateserver

import (
	"github.com/mx5566/server/entity"
	"reflect"
	"sync"
)

var GSessionMgr = CreateSessionMgr()

type ISessionMgr interface {
}

type SessionMgr struct {
	entity.Entity
	entity.EntityPool

	sync.RWMutex
	sessions map[int64]entity.IEntity
}

func CreateSessionMgr() *SessionMgr {
	s := &SessionMgr{}

	s.Init()

	return s
}

func (m *SessionMgr) Init() {
	m.sessions = make(map[int64]entity.IEntity)

	m.Entity.Init()
	m.EntityPool.InitPool(reflect.TypeOf(ClientSession{}))

	entity.RegisterEntity(m)
}

func (m *SessionMgr) Update() {
	temp := make(map[int64]entity.IEntity)
	m.RWMutex.Lock()
	for k, v := range m.sessions {
		temp[k] = v
	}
	m.sessions = make(map[int64]entity.IEntity)
	m.RWMutex.Unlock()

	for _, v := range temp {
		m.EntityPool.AddEntity(v)
	}

	m.EntityPool.Update()
}

func (m *SessionMgr) AddSession(s entity.IEntity) {
	m.RWMutex.Lock()
	m.sessions[s.GetID()] = s
	m.RWMutex.Unlock()

	//m.EntityPool.AddEntity(s)
}

func (m *SessionMgr) DelSession(ID int64) {
	m.EntityPool.DelEntity(ID)
}
