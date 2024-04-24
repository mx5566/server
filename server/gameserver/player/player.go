package player

import "github.com/mx5566/server/base/entity"

type Player struct {
	entity.Entity
}

func (p *Player) Init() {
	p.Entity.Init()
	p.Entity.Start()
}
