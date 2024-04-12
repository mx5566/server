package aoi

import (
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

type Player struct {
	aoi AOI
}

func (p *Player) OnEnterAoi(aoi *AOI) {

}

func (p *Player) OnLeaveAoi(aoi *AOI) {
}

func randPos(min, max int32) int32 {
	return min + rand.Int31n(max-min)
}

func TestAoi(t *testing.T) {
	mgr := NewAoiAmager(100)

	players := []*Player{}
	for i := 0; i < 1000; i++ {
		player := &Player{}
		InitAoi(&player.aoi, player)
		players = append(players, player)
		mgr.Enter(&player.aoi, randPos(10, 500), randPos(10, 500))
	}

	proffd, _ := os.OpenFile("test_aoi"+".pprof", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer proffd.Close()

	pprof.StartCPUProfile(proffd)
	for i := 0; i < 10; i++ {
		t0 := time.Now()
		for _, obj := range players {
			mgr.Move(&obj.aoi, obj.aoi.x+randPos(-10, 10), obj.aoi.y+randPos(-10, 10))
			mgr.Leave(&obj.aoi)
			mgr.Enter(&obj.aoi, obj.aoi.x+randPos(-10, 10), obj.aoi.y+randPos(-10, 10))
		}
		dt := time.Now().Sub(t0)
		t.Logf("%d objects takes %s", 1000, dt)
	}

	for _, obj := range players {
		mgr.Leave(&obj.aoi)
	}
	pprof.StopCPUProfile()

}
