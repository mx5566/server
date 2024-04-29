package aoiskiplist

import (
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/entity"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func (p *Player) OnLeaveAoi(aoi *AOIProc) {
	/*log.Println("OnLeaveAoi-------------------", aoi.ID, "---", p.aoi.ID)
	 */
}

func (p *Player) OnEnterAoi(aoi *AOIProc) {
	/*log.Println("OnEnterAoi-------------------", aoi.ID, "---", p.aoi.ID)
	 */
	//entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, "Player.interest", aoi.unit)
}

type Player struct {
	entity.Entity
	aoi AOIProc
}

func TestMap(t *testing.T) {
	mgr := NewAoiAmager(5)

	players := []*Player{}
	for i := 0; i < 1000; i++ {
		player := &Player{}
		player.aoi.ID = int64(10 + i)
		InitAoi(&player.aoi, player, player)

		players = append(players, player)
		mgr.Enter(&player.aoi, player.aoi.x+10.0, player.aoi.y+10.0)

		//t.Logf("mark1:  %d\n", player.aoi.mark)
	}

	proffd, _ := os.OpenFile("test_aoi_proc"+".pprof", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer proffd.Close()

	pprof.StartCPUProfile(proffd)
	for i := 0; i < 10; i++ {
		t0 := time.Now()
		for _, obj := range players {
			mgr.Move(&obj.aoi, obj.aoi.x+10.0, obj.aoi.y+10.0)
			mgr.Leave(&obj.aoi)
			mgr.Enter(&obj.aoi, obj.aoi.x+10.0, obj.aoi.y+10.0)
		}
		dt := time.Now().Sub(t0)
		t.Logf("%d objects takes %s", 1000, dt)
	}

	for _, obj := range players {
		mgr.Leave(&obj.aoi)
	}

	pprof.StopCPUProfile()
}

func randPos(min, max int32) int32 {
	return min + rand.Int31n(max-min)
}
func TestAoi(t *testing.T) {
	mgr := NewAoiAmager(5)

	//time.Sleep(20 * time.Second)
	t0 := time.Now()

	players := []*Player{}
	for i := 0; i < 1000; i++ {
		player := &Player{}
		player.aoi.ID = int64(10 + i)
		InitAoi(&player.aoi, player, player)

		players = append(players, player)
		mgr.Enter(&player.aoi, base.Coord(randPos(10, 500)), base.Coord(randPos(10, 500)))

		//t.Logf("mark1:  %d\n", player.aoi.mark)
	}

	dt := time.Now().Sub(t0)
	t.Logf("Enter time:%s", dt)

	proffd, _ := os.OpenFile("test_aoi_proc"+".pprof", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer proffd.Close()

	//trace.Start(os.Stderr)

	f, _ := os.Create("mem.pprof")
	defer f.Close()
	pprof.StartCPUProfile(proffd)
	for i := 0; i < 1000; i++ {
		t0 := time.Now()
		for _, obj := range players {
			mgr.Move(&obj.aoi, obj.aoi.x+base.Coord(randPos(-10, 10)), obj.aoi.y+base.Coord(randPos(-10, 10)))

			//t.Logf("mark2:  %d\n", obj.aoi.mark)

			mgr.Leave(&obj.aoi)
			mgr.Enter(&obj.aoi, obj.aoi.x+base.Coord(base.Random[float32](-10, 10)), obj.aoi.y+base.Coord(base.Random[float32](-10, 10)))
		}
		dt := time.Now().Sub(t0)
		t.Logf("%d objects takes %s count:%d ", 1000, dt, i)
	}

	for _, obj := range players {
		mgr.Leave(&obj.aoi)
	}

	pprof.StopCPUProfile()

	pprof.WriteHeapProfile(f)

	//trace.Start(os.Stderr)

}
