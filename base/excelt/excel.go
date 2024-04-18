package excelt

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/mx5566/logm"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 加载table
func init() {
	//Load()
}

func LoadExcel(path string) {

}

func compressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

// 把key转换位字符串
func CombineKeys(keys ...interface{}) string {
	//sort.Strings(keys)
	fmt.Println(keys...)
	com := []string{}
	for _, key := range keys {
		switch key.(type) {
		case int, int32, int64, int8, int16:
			com = append(com, strconv.FormatInt(reflect.ValueOf(key).Int(), 10))
		case uint, uint32, uint64, uint16, uint8:
			com = append(com, strconv.FormatUint(reflect.ValueOf(key).Uint(), 10))
		case string:
			com = append(com, key.(string))
		default:
			fmt.Println("unkonw type "+reflect.TypeOf(key).String(), " ", key)
		}
	}
	return strings.Join(com, "_")
}

// 把key转换位字符串
func CombineKeysEx(keys []interface{}) string {
	//sort.Strings(keys)
	com := []string{}
	for _, key := range keys {
		switch key.(type) {
		case int, int32, int64, int8, int16:
			com = append(com, strconv.FormatInt(reflect.ValueOf(key).Int(), 10))
		case uint, uint32, uint64, uint16, uint8:
			com = append(com, strconv.FormatUint(reflect.ValueOf(key).Uint(), 10))
		case string:
			com = append(com, key.(string))
		default:
			logm.ErrorfE("unkonw type %s  key %v", reflect.TypeOf(key).String(), key)
		}
	}
	return strings.Join(com, "_")
}

func Read(fileName string, keys ...string) map[interface{}][]byte {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		println(err.Error())
		return nil
	}

	// 找出key
	mapKeys := make(map[string]interface{})
	for _, value := range keys {
		fmt.Println(value)
		mapKeys[value] = value
	}

	var mapFields = make(map[interface{}]map[string]interface{})
	var mapFieldsBytes = make(map[interface{}][]byte)
	var mapFieldNames = make(map[string]string)
	var sliceFieldNames = []string{}
	var sliceFieldTypes = []string{}
	// 获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	for index, row := range rows {
		// 第一行算是一种注释
		if index == 0 {
			for _, colCell := range row {
				if colCell == "" {
					log.Panic("fileName " + fileName + " has field empty 0!!!")
				}
				//fmt.Print(colCell)
				mapFieldNames[colCell] = colCell
			}
			continue
		}

		// 第二行是字段名字
		if index == 1 {
			for _, colCell := range row {
				if colCell == "" {
					log.Panic("fileName " + fileName + " has field empty 1!!!")
				}
				//fmt.Print(colCell)
				sliceFieldNames = append(sliceFieldNames, colCell)
			}
			continue
		}

		// 第三行是数据类型
		if index == 2 {
			for _, colCell := range row {
				if colCell == "" {
					log.Panic("fileName " + fileName + " has field empty 2!!!")
				}

				colCell = compressStr(colCell)
				//fmt.Print(colCell)
				sliceFieldTypes = append(sliceFieldTypes, colCell)
			}
			continue
		}

		oneMapFields := make(map[string]interface{})
		oneMapFieldsBytes := []byte{}
		comKeys := []string{}
		for index1, colCell := range row {
			// 实际的值判断
			fieldName := sliceFieldNames[index1]
			if _, ok := mapKeys[fieldName]; ok {
				comKeys = append(comKeys, colCell)
			}

			switch sliceFieldTypes[index1] {
			case "int64", "int32", "int":
				ret, _ := strconv.Atoi(colCell)
				oneMapFields[fieldName] = ret
			case "float32":
				//ret, _ := strconv.Atoi(colCell)
				//strconv.FormatFloat(float64, 'E', -1, 32)
				ret, _ := strconv.ParseFloat(colCell, 32)
				oneMapFields[fieldName] = float32(ret)
			case "float64":
				ret, _ := strconv.ParseFloat(colCell, 64)
				oneMapFields[fieldName] = ret
			case "string":
				oneMapFields[fieldName] = colCell
			case "[]int":
				sli := strings.Split(colCell, ",")
				sliTemp := []int{}
				for _, value := range sli {
					ret, _ := strconv.Atoi(value)
					sliTemp = append(sliTemp, ret)
				}
				// 设置数组
				oneMapFields[fieldName] = sliTemp
			case "[]string":
				sli := strings.Split(colCell, "|")
				// 设置数组
				oneMapFields[fieldName] = sli
			case "map[string]string": // key1,value1|key2,value2

			}
		}
		//sort.Strings(comKeys)
		oneMapFieldsBytes, err = json.Marshal(oneMapFields)
		if err != nil {
			log.Panic("json.Marshal table fileName error ", err)
		}
		mapFields[strings.Join(comKeys, "_")] = oneMapFields
		mapFieldsBytes[strings.Join(comKeys, "_")] = oneMapFieldsBytes
	}

	return mapFieldsBytes
}

func ListFileFunc(p []string) {
	for index, value := range p {
		fmt.Println("Index = ", index, " Value = ", value)
		if index == 0 {
			Read(value, "ID")
		}
	}
}
