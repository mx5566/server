package network

import "encoding/binary"

type MsgPacket struct {
	MsgId   uint32
	MsgBody []byte
}

type DataPacket struct {
}

func (d DataPacket) Encode(msg *MsgPacket) []byte {
	dataLen := len(msg.MsgBody)
	buff := make([]byte, TcpHeadSize+dataLen)
	binary.BigEndian.PutUint32(buff[0:TcpIDLength], msg.MsgId)
	binary.BigEndian.PutUint16(buff[TcpIDLength:TcpHeadSize], uint16(dataLen))
	copy(buff[TcpHeadSize:], msg.MsgBody)

	return buff
}

func (d DataPacket) Decode(buff []byte) *MsgPacket {
	msgId := binary.BigEndian.Uint32(buff[0:TcpIDLength])
	msgLen := binary.BigEndian.Uint16(buff[TcpIDLength:TcpHeadSize])
	msgBody := buff[TcpHeadSize : TcpHeadSize+msgLen]

	msg := new(MsgPacket)
	msg.MsgId = msgId
	msg.MsgBody = msgBody

	return msg
}
