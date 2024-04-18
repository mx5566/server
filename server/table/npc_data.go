//go:noformat
package table

// Auto generator code, Do not edit it.
import (
	"encoding/json"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/excelt"
)

//
type NpcBase struct {
	Type int
	Level int
	Hp int64
	AttackInter int32
	AttackDistance float32
	ID int
	Name string
}

var MapNpc map[interface{}]NpcBase

func LoadNpcTable(path string) {
	data := excelt.Read(path)

	MapNpc = make(map[interface{}]NpcBase)
	for key, value := range data {
		var tableBase NpcBase
		err := json.Unmarshal(value, &tableBase)
		if err != nil {
			logm.DebugfE("load Npc table err key:%s,  error:%s\n", key, err.Error())
			continue
		}

		MapNpc[key] = tableBase
	}
}

func GetNpcBase(keys ...interface{}) *NpcBase {
	keyCom := excelt.CombineKeysEx(keys)
	if base, ok := MapNpc[keyCom]; ok {
		return &base
	}

	return nil
}
