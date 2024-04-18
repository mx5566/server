//go:generate  go run generate.go
package main

import (
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/excelt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	generate("../../release/table/")
}

func generate(path string) {
	var filter base.FileFilter
	_ = filter.GetFileList(path, base.Xlsx)

	list := filter.ListFile

	//logm.DebugfE("xlsx file load count[%v], list[%v]", len(list), list)

	// 把所有的文件获取文件名字
	for k, v := range list {
		b := filepath.Base(v)
		list[k] = b
	}

	for _, file := range list {
		temPath := path + file
		header := excelt.ReadBase(temPath)
		generateTableCode(header)
	}
}

func generateTableCode(header *excelt.TableHeader) {
	upperName := base.Capitalize(header.TableName)

	// 创建类
	class := "type " + upperName + "Base struct {\n"
	for key, value := range header.FieldNameType {
		class += fmt.Sprintf("	%s %s\n", key, value)
	}

	class += "}"

	tempCode := templateCode

	tempCode = strings.Replace(tempCode, "{class}", class, -1)
	tempCode = strings.Replace(tempCode, "{tableName}", upperName, -1)

	f, _ := os.Create("../../server/table/" + header.TableName + "_data.go")
	f.WriteString(tempCode)
	f.Sync()
	f.Close()

	logm.DebugfE("生成表:%s, 名字: %s", header.TableName, header.TableName+"_data.go")
}

// 自动去生成对应表的go文件  模版
var templateCode = `//go:noformat
package table

// Auto generator code, Do not edit it.
import (
	"encoding/json"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/excelt"
)

//
{class}

var Map{tableName} map[interface{}]{tableName}Base

func Load{tableName}Table(path string) {
	data := excelt.Read(path)

	Map{tableName} = make(map[interface{}]{tableName}Base)
	for key, value := range data {
		var tableBase {tableName}Base
		err := json.Unmarshal(value, &tableBase)
		if err != nil {
			logm.DebugfE("load {tableName} table err key:%s,  error:%s\n", key, err.Error())
			continue
		}

		Map{tableName}[key] = tableBase
	}
}

func Get{tableName}Base(keys ...interface{}) *{tableName}Base {
	keyCom := excelt.CombineKeysEx(keys)
	if base, ok := Map{tableName}[keyCom]; ok {
		return &base
	}

	return nil
}
`
