package excel

import (
	"encoding/json"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/excelt"
	"path/filepath"
	"runtime"
)

// TODO 查找所有的TODO位置
var (
	ostype = runtime.GOOS // 获取系统类型
)

// 所有表的预定义字符串
var (
	ItemTableStr  = "item"
	EquipTableStr = "equip"
	NpcTableStr   = "npc"
	SkillTableStr = "skill"
)

func Load(path string) {
	var filter base.FileFilter
	_ = filter.GetFileList(path, base.Xlsx)

	list := filter.ListFile

	//logm.DebugfE("xlsx file load count[%v], list[%v]", len(list), list)

	// 把所有的文件获取文件名字
	for k, v := range list {
		b := filepath.Base(v)
		list[k] = b
	}
	// TODO 按照示例添加表
	for _, file := range list {
		temPath := path + file
		switch file {
		case "equip.xlsx":
			LoadEquipTable(temPath)
		case "item.xlsx":
			LoadItemTable(temPath)
		case "npc.xlsx":
			LoadNpcTable(temPath)
		default:
			fmt.Println("error path " + file)
		}
	}
}

var MapItemsBase map[interface{}]ItemBase
var MapNpcBase map[interface{}]NpcBase
var MapEquipsBase map[interface{}]EquipBase

type ItemBase struct {
	ID       int64    `json:"ID"`
	Name     string   `json:"Name"`
	Type     uint16   `json:"Type"`
	Quality  uint8    `json:"Quality"`
	Ratio1   float32  `json:"Ratio1"`
	Ratio2   float64  `json:"Ratio2"`
	BufferID []int32  `json:"BufferID"`
	Names    []string `json:"Names"`
}

type NpcBase struct {
	ID             int64   `json:"ID"`
	Name           string  `json:"Name"`
	Type           uint16  `json:"Type"`
	Level          uint16  `json:"Level"`
	Hp             int64   `json:"Hp"`
	AttackInter    int32   `json:"AttackInter"`
	AttackDistance float32 `json:"AttackDistance"`
}

type EquipBase struct {
	ItemBase
	// external attr
}

// //////////////////////////////////////////////////////////////////////////
func LoadItemTable(path string) {
	items := excelt.Read(path, "ID")

	fmt.Println("load table item !!!")
	//fmt.Println(items)

	MapItemsBase = make(map[interface{}]ItemBase)
	for key, value := range items {
		var itemBase ItemBase
		err := json.Unmarshal(value, &itemBase)
		if err != nil {
			fmt.Println("load item table LoadItem err key [ ", key, "]  error [", err, " ]")
			continue
		}
		MapItemsBase[key] = itemBase
	}
}

func LoadEquipTable(path string) {
	_ = excelt.Read(path, "ID")
	fmt.Println("load table equip !!!")

}

func LoadNpcTable(path string) {
	npcs := excelt.Read(path, "ID")

	fmt.Println("load table item !!!")
	//fmt.Println(items)

	MapNpcBase = make(map[interface{}]NpcBase)
	for key, value := range npcs {
		var npcBase NpcBase
		err := json.Unmarshal(value, &npcBase)
		if err != nil {
			fmt.Println("load item table LoadItem err key [ ", key, "]  error [", err, " ]")
			continue
		}
		MapNpcBase[key] = npcBase
	}
}

// TODO: 需要加入对用的表返回
func GetBase(name string, keys ...interface{}) interface{} {
	if len(name) == 0 {
		return nil
	}

	logm.DebugfE("GetBase name[%s] key[%v]", name, keys)
	keyCom := excelt.CombineKeysEx(keys)
	switch name {
	case ItemTableStr:
		if base, ok := MapItemsBase[keyCom]; ok {
			return &base
		}
	case EquipTableStr:
		if base, ok := MapEquipsBase[keyCom]; ok {
			return &base
		}
	case NpcTableStr:
		if base, ok := MapNpcBase[keyCom]; ok {
			return &base
		}
	default:
		logm.ErrorfE("GetBase Error name %s", name)
	}

	return nil
}
