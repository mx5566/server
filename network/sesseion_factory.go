package network

type SESSION_TYPE int

const (
	SESSION_NON SESSION_TYPE = iota
	SESSION_CLIENT
	SESSION_GAME
	SESSION_WORLD
	SESSION_LOGIN
	SESSION_GATE
	SESSION_MAX
)

type ISessionFactory interface {
	CreateSession() ISession
	InitSessionPacket()
	GetID() int
}

var GSessionFactoryMgr SessionFactoryMgr

type SessionFactoryMgr struct {
	factorys map[int]ISessionFactory
}

func CreateFactoryMgr() *SessionFactoryMgr {
	f := &SessionFactoryMgr{
		factorys: make(map[int]ISessionFactory),
	}

	return f
}

func (s *SessionFactoryMgr) GetFactory(Id int /*类型*/) ISessionFactory {
	if _, ok := s.factorys[Id]; ok {
		return s.factorys[Id]
	}

	return nil
}

func (s *SessionFactoryMgr) AddFactory(factory ISessionFactory) {
	s.factorys[factory.GetID()] = factory
}
