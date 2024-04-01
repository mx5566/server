package entity

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/server/pb"
)

var GEntityMgr EntityMgr

type EntityMgr struct {
	Entitys map[string]IEntity // 全局注册 一个类只有个注册，不管有多少个对象
}

// 注册是全局的 只有一次注册
// 对于需要多次注册的实体 比如玩家作为实体注册，我们不会把全部的玩家都注册进来因为我们调用其实只是获取实体类的方法，而不是单个对象的方法
// 但是我们有需要保存所有的实体对象
func (m *EntityMgr) RegisterEntity(entity IEntity) {
	name := base.GetClassName(entity)

	if _, ok := m.Entitys[name]; !ok {
		logm.PanicfE("重复的注册实体Name: %s \n", name)
		return
	}

	// 记录实体名字和实体的映射
	m.Entitys[name] = entity
}

func (m *EntityMgr) Call(packet pb.RpcPacket) {
	className := packet.Head.ClassName
	funcName := packet.Head.FuncName
	if v, ok := m.Entitys[className]; ok && v != nil {
		if v.IsExistMethod(funcName) {
			v.Call(packet)
		}
	}

}

func RegisterEntity(entity IEntity) {
	GEntityMgr.RegisterEntity(entity)
}
