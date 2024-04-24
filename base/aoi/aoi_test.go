package aoi

import (
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func (p *Player) OnLeaveAoi(aoi *AOI) {
}

func (p *Player) OnEnterAoi(aoi *AOI) {
	entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, "Player.interest", aoi.unit)
}

func randPos(min, max int32) int32 {
	return min + rand.Int31n(max-min)
}

type FriendData struct {
	ID          int64
	FriendValue int64
}

type Player struct {
	entity.Entity
	aoi AOI

	// 自己的数据
	// 好有数据
	Frinds map[int64]FriendData
}

func TestMap(t *testing.T) {

}

func TestAoi(t *testing.T) {
	mgr := NewAoiAmager(100)

	players := []*Player{}
	for i := 0; i < 10; i++ {
		player := &Player{}
		InitAoi(&player.aoi, player, player)
		players = append(players, player)
		mgr.Enter(&player.aoi, float32(randPos(10, 500)), float32(randPos(10, 500)))
	}

	proffd, _ := os.OpenFile("test_aoi"+".pprof", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer proffd.Close()

	pprof.StartCPUProfile(proffd)
	for i := 0; i < 10; i++ {
		t0 := time.Now()
		for _, obj := range players {
			mgr.Move(&obj.aoi, obj.aoi.x+float32(randPos(-10, 10)), obj.aoi.y+float32(randPos(-10, 10)))
			mgr.Leave(&obj.aoi)
			mgr.Enter(&obj.aoi, obj.aoi.x+float32(randPos(-10, 10)), obj.aoi.y+float32(randPos(-10, 10)))
		}
		dt := time.Now().Sub(t0)
		t.Logf("%d objects takes %s", 1000, dt)
	}

	for _, obj := range players {
		mgr.Leave(&obj.aoi)
	}
	pprof.StopCPUProfile()

}
