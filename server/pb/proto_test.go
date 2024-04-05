package pb

import (
	"fmt"
	"github.com/mx5566/server/rpc3"
	"testing"
)

func TestA(t *testing.T) {

	//var t1 = &pb
	//name := proto.MessageReflect(t1).Descriptor().FullName()

	//t.Logf("name: %s\n", name)

	head := rpc3.RpcHead{}
	funcName := ""
	cluster := rpc3.ClusterInfo{
		Ip:          "0.0.0.0",
		Port:        8080,
		ServiceType: 1,
	}

	rpcPacket := Marshal(&head, &funcName, &cluster)

	fmt.Println(rpcPacket)
}
