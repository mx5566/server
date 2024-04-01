package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/network"
	"github.com/mx5566/server/server/pb"
	"hash/crc32"
	"time"
)

type clientData struct {
	ID int32
}

func main() {

	client := new(network.ClientSocket)
	client.Init("127.0.0.1", 8080)
	client.Start()

	data := pb.Test{
		Name:     "mengxiang",
		PassWord: "990000",
	}

	serData, _ := proto.Marshal(&data)
	fmt.Println("send data len ", len(serData))

	msg := new(network.MsgPacket)
	msg.MsgId = crc32.ChecksumIEEE([]byte(base.GetMessageName(&data)))
	msg.MsgBody = serData

	dp := network.DataPacket{}
	buff := dp.Encode(msg)

	client.Send(buff)

	time.Sleep(100 * time.Second)

	client.Stop()
}
