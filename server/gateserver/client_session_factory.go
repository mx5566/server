package gateserver

import "github.com/mx5566/server/network"

type ClientSessionFactory struct {
}

func (f *ClientSessionFactory) CreateSession() network.ISession {
	s := NewSession()
	return s
}

func (f *ClientSessionFactory) InitSessionPacket() {

}

func (f *ClientSessionFactory) GetID() int {
	return int(network.SESSION_CLIENT)
}
