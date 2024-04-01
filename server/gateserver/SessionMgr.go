package gateserver

import (
	"github.com/mx5566/server/entity"
	"github.com/mx5566/server/network"
	"reflect"
)

var GSessionMgr = CreateSessionMgr()

type ISessionMgr interface {
}

type SessionMgr struct {
	entity.Entity
	entity.EntityPool

	sessions map[int64]network.ISession
}

func CreateSessionMgr() *SessionMgr {
	s := &SessionMgr{}

	s.Init()

	return s
}

func (m *SessionMgr) Init() {
	m.Entity.Init()
	m.EntityPool.InitPool(reflect.TypeOf(ClientSession{}))

	entity.RegisterEntity(m)
}

func (m *SessionMgr) Update() {
	m.EntityPool.Update()
}

func (m *SessionMgr) AddSession(s entity.IEntity) {
	m.EntityPool.AddEntity(s)
}

func (m *SessionMgr) DelSession(ID int64) {
	m.EntityPool.DelEntity(ID)
}
