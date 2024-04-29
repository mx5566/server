package player

import (
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/aoi"
	"github.com/mx5566/server/base/entity"
	"reflect"
)

type Player struct {
	entity.Entity
	Unit
}

func (p *Player) Init() {
	p.Entity.Init()
	p.Entity.Start()
	p.Unit.Init()
	p.TypeName = reflect.TypeOf(p).Name()
	aoi.InitAoi(&p.Aoi, &p.Unit, &p.Unit)
}

func (p *Player) SetPosition(x, y, z base.Coord) {
	p.Unit.Vector3 = base.Vector3{x, y, z}
}

func (p *Player) SetMapID(mapID uint32) {
	p.Unit.MapID = mapID
}
