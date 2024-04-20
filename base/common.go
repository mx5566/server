package base

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"hash/crc32"
	"reflect"
	"runtime"
	"strings"
)

const (
	Send_Game = iota
	Send_Gate
	Send_Login
)

// 服务器注册的前缀
const ServiceName = "Server/"

// 模块注册的前缀
const ModuleNameDir = "Module/"

// 对象注册的前缀
const ObjectDir = "Object/"

const (
	Json = ".json"
	Xlsx = ".xlsx"
)

type LoginState uint8

const (
	LoginState_None            LoginState = iota
	LoginState_AccountLogining            // 账号登录中
	LoginState_AccountLogin               // 账号登陆了 进入角色选择界面
	LoginState_PlayerLogin                // 玩家选择角色登录游戏了
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
	logm.ErrorfE("==> %s", data)
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

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
