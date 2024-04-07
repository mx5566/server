package base

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"hash/crc32"
	"log"
	"reflect"
	"runtime"
	"strings"
)

const (
	Send_Game = iota
	Send_Gate
	Send_Login
)

const ServiceName = "Server/"

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

func GetProMessageByName(name string) proto.Message {
	messageType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(name))
	if err != nil {
		logm.ErrorfE("根据proto消息的名字获取类型失败: %s err : %s", name, err.Error())
		return nil
	}
	msg := proto.MessageV1(messageType.New())

	return msg
}

func GetMessageID(msg proto.Message) uint32 {
	name := GetMessageName(msg)
	return crc32.ChecksumIEEE([]byte(name))
}

func GetClassName(class interface{}) string {
	rType := reflect.TypeOf(class)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	return rType.Name()
}
