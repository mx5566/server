package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/pb"
	"hash/crc32"
	"time"
)

type clientData struct {
	ID int32
}

func HandleMsg(packet rpc3.Packet) {
	logm.DebugfE("接收到数据了")
}

func main() {

	//var id int64 = 0

	clients := make(map[int]network.ISocket)
	for i := 0; i < 1; i++ {
		client := new(network.ClientSocket)
		client.Init("127.0.0.1", 8080)
		client.Start()
		clients[i] = client
		client.BindPacketFunc(HandleMsg)
		//client.GetConnId()
		//ii := atomic.AddInt64(&id, 1)
		data := pb.LoginAccountReq{
			UserName: "mengxiang",
			Password: "9090",
		}

		serData, _ := proto.Marshal(&data)
		fmt.Println("send data len ", len(serData), "i :", i)

		msg := new(network.MsgPacket)
		msg.MsgId = crc32.ChecksumIEEE([]byte(base.GetMessageName(&data)))
		msg.MsgBody = serData

		dp := network.DataPacket{}
		buff := dp.Encode(msg)

		client.Send(rpc3.Packet{Buff: buff})
	}

	time.Sleep(100 * time.Second)

	for _, v := range clients {
		v.Stop()
	}
}
