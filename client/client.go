package main

import (
	"encoding/json"
	"fmt"
	"github.com/mx5566/server/network"
	"time"
)

type clientData struct {
	ID int32
}

func main() {
	client := new(network.ClientSocket)
	client.Init("127.0.0.1", 8080)
	client.Start()

	data := clientData{
		ID: 10,
	}

	serData, _ := json.Marshal(&data)
	fmt.Println("send data len ", len(serData))

	msg := new(network.MsgPacket)
	msg.MsgId = 1001
	msg.MsgBody = serData

	dp := network.DataPacket{}
	buff := dp.Encode(msg)

	client.Send(buff)

	time.Sleep(10 * time.Second)

	client.Stop()
}
