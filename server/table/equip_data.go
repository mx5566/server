//go:noformat
package table

// Auto generator code, Do not edit it.
import (
	"encoding/json"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/excelt"
)

//
type EquipBase struct {
	Ratio1 float32
	Ratio2 float64
	Ids []int
	Names []string
	ID int
	Name string
	Type int
	Quality int
}

var MapEquip map[interface{}]EquipBase

func LoadEquipTable(path string) {
	data := excelt.Read(path)

	MapEquip = make(map[interface{}]EquipBase)
	for key, value := range data {
		var tableBase EquipBase
		err := json.Unmarshal(value, &tableBase)
		if err != nil {
			logm.DebugfE("load Equip table err key:%s,  error:%s\n", key, err.Error())
			continue
		}

		MapEquip[key] = tableBase
	}
}

func GetEquipBase(keys ...interface{}) *EquipBase {
	keyCom := excelt.CombineKeysEx(keys)
	if base, ok := MapEquip[keyCom]; ok {
		return &base
	}

	return nil
}
