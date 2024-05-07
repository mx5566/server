package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/pb"
	"hash/crc32"
	"log"
	"sync/atomic"
	"time"
)

type clientData struct {
	ID int32
}

type Func func(data []byte)

var HandlePackets map[uint32]Func

func Register(msg proto.Message, f Func) {
	id := base.GetMessageID(msg)

	HandlePackets[id] = f
}

func HandlAccountRep(data []byte) {
	pMsg := pb.LoginAccontRep{}

	err := proto.Unmarshal(data, &pMsg)
	if err != nil {
		log.Println(err)
		return
	}

	if pMsg.ErrCode != 0 {
		log.Printf("账号登录失败: %d\n", pMsg.ErrCode)
		return
	}

	account := new(Account)
	account.ID = pMsg.AccountId
	for i := 0; i < len(pMsg.PList); i++ {
		pl := new(Player)
		pl.Name = pMsg.PList[i].PlayerName
		pl.ID = pMsg.PList[i].PlayerId
		account.players = append(account.players, pl)
	}
	accountS = account

	if len(pMsg.PList) >= 3 {
		PlayerLoginRequest(clients[0])
	} else {
		CreatePlayerRequest(clients[0])
	}

	log.Printf("accountLogin code:%d, list:%v, accountId:%d\n", pMsg.ErrCode, pMsg.PList, pMsg.AccountId)
}

func HandlePlayerLogin(data []byte) {
	pMsg := pb.LoginPlayerRep{}

	err := proto.Unmarshal(data, &pMsg)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("playerLogin code:%d, playerID:%d\n", pMsg.ErrCode, pMsg.PlayerId)

}

func HandleCreatePlayer(data []byte) {
	pMsg := pb.CreatePlayerRep{}

	err := proto.Unmarshal(data, &pMsg)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("createPlayer code:%d, playerID:%d, name:%s\n", pMsg.ErrCode, pMsg.PlayerId, pMsg.Name)
}

func HandleRoleList(data []byte) {
	pMsg := pb.RoleSelectListRep{}

	err := proto.Unmarshal(data, &pMsg)
	if err != nil {
		log.Println(err)
		return
	}

	accountS.players = make([]*Player, 0)
	for i := 0; i < len(pMsg.PList); i++ {
		pl := new(Player)
		pl.Name = pMsg.PList[i].PlayerName
		pl.ID = pMsg.PList[i].PlayerId

		accountS.players = append(accountS.players, pl)
	}

	log.Printf("roleList accountId:%d, list:%v\n", pMsg.AccountId, pMsg.GetPList())
}

func init() {
	HandlePackets = make(map[uint32]Func)
	Register(&pb.LoginAccontRep{}, HandlAccountRep)
	Register(&pb.LoginPlayerRep{}, HandlePlayerLogin)
	Register(&pb.CreatePlayerRep{}, HandleCreatePlayer)
	Register(&pb.RoleSelectListRep{}, HandleRoleList)
}

func HandleMsg(packet rpc3.Packet) {
	logm.DebugfE("接收到数据了")

	buff := packet.Buff

	dp := network.DataPacket{}
	msg := dp.Decode(buff)

	_, ok := HandlePackets[msg.MsgId]
	if !ok {
		return
	}

	HandlePackets[msg.MsgId](msg.MsgBody)
}

var ss atomic.Int32

func CreatePlayerRequest(client network.ISocket) {
	data := pb.CreatePlayerReq{
		Name:      "role1" + string(base.RandomInt(1, 1000)),
		AccountID: accountS.ID,
	}

	serData, _ := proto.Marshal(&data)

	msg := new(network.MsgPacket)
	msg.MsgId = crc32.ChecksumIEEE([]byte(base.GetMessageName(&data)))
	msg.MsgBody = serData

	dp := network.DataPacket{}
	buff := dp.Encode(msg)

	client.Send(rpc3.Packet{Buff: buff})
}

func AccountLoginRequest(client network.ISocket) {
	data := pb.LoginAccountReq{
		UserName: "mengxiang",
		Password: "9090",
	}

	serData, _ := proto.Marshal(&data)

	msg := new(network.MsgPacket)
	msg.MsgId = crc32.ChecksumIEEE([]byte(base.GetMessageName(&data)))
	msg.MsgBody = serData

	dp := network.DataPacket{}
	buff := dp.Encode(msg)

	client.Send(rpc3.Packet{Buff: buff})

	log.Printf("发送账号登录的请求")
}

func PlayerLoginRequest(client network.ISocket) {
	data := pb.LoginPlayerReq{
		PlayerId:  accountS.players[0].ID,
		AccountID: accountS.ID,
	}

	serData, _ := proto.Marshal(&data)

	msg := new(network.MsgPacket)
	msg.MsgId = crc32.ChecksumIEEE([]byte(base.GetMessageName(&data)))
	msg.MsgBody = serData

	dp := network.DataPacket{}
	buff := dp.Encode(msg)

	client.Send(rpc3.Packet{Buff: buff})

	log.Printf("发送角色登录请求")
}

var accountS *Account

type Account struct {
	ID      int64
	players []*Player
}

type Player struct {
	Name string
	ID   int64
}

var clients map[uint32]network.ISocket

func main() {

	//var id int64 = 0

	clients = make(map[uint32]network.ISocket)
	for i := 0; i < 1; i++ {
		client := new(network.ClientSocket)
		client.Init("127.0.0.1", 13000)
		client.Start()
		clients[client.GetConnId()] = client
		client.BindPacketFunc(HandleMsg)
		//client.GetConnId()
		//ii := atomic.AddInt64(&id, 1)

		AccountLoginRequest(client)

		time.Sleep(3 * time.Second)

	}

	time.Sleep(100 * time.Second)

	for _, v := range clients {
		v.Stop()
	}
}
