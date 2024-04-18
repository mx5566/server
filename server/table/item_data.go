//go:noformat
package table

// Auto generator code, Do not edit it.
import (
	"encoding/json"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/excelt"
)

//
type ItemBase struct {
	Type int
	Quality int
	Ratio1 float32
	Ratio2 float64
	Ids []int
	Names []string
	ID int
	Name string
}

var MapItem map[interface{}]ItemBase

func LoadItemTable(path string) {
	data := excelt.Read(path)

	MapItem = make(map[interface{}]ItemBase)
	for key, value := range data {
		var tableBase ItemBase
		err := json.Unmarshal(value, &tableBase)
		if err != nil {
			logm.DebugfE("load Item table err key:%s,  error:%s\n", key, err.Error())
			continue
		}

		MapItem[key] = tableBase
	}
}

func GetItemBase(keys ...interface{}) *ItemBase {
	keyCom := excelt.CombineKeysEx(keys)
	if base, ok := MapItem[keyCom]; ok {
		return &base
	}

	return nil
}
