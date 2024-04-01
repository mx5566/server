package base

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
	"runtime"
	"strings"
)

// 输出错误，跟踪代码
func TraceCode(code ...interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	data := ""
	for _, v := range code {
		data += fmt.Sprintf("%v", v)
	}
	data += string(buf[:n])
	log.Printf("==> %s\n", data)
}

func GetMessageName(msg proto.Message) string {
	name := proto.MessageReflect(msg).Descriptor().FullName()
	if i := strings.LastIndexByte(string(name), '.'); i >= 0 {
		name = name[i+1:]
	}

	return string(name)
}

func GetClassName(class interface{}) string {
	rType := reflect.TypeOf(class)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	return rType.Name()
}
