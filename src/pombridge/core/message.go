package core

import "encoding/binary"

/*
	Message Packet:

	| 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
	|-------------------------------|
0 |   packet len  |version|syn|fin|
  |-------------------------------|
8 |      seq      |      ack      |
	|-------------------------------|
16|    channel    |    data..     |
	|-------------------------------|
*/

type Message struct {
	flow     *Flow
	syn, fin bool
	seq, ack uint16
	channel  uint16
	data     []byte
}

type MsgChan chan *Message

const (
	versionMajor    = 1
	versionMinor    = 0
	PacketHeaderLen = 20
)

func boolToByte(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}

func bytesToUint16(b []byte) uint16 {
	val, _ := binary.Uvarint(b)
	return uint16(val)
}

func (msg *Message) PacketHeader(buf []byte) {
	packetLen := PacketHeaderLen + len(msg.data)

	binary.PutUvarint(buf[0:4], uint64(packetLen))
	buf[4] = versionMajor
	buf[5] = versionMinor
	buf[6] = boolToByte(msg.syn)
	buf[7] = boolToByte(msg.fin)
	binary.PutUvarint(buf[8:12], uint64(msg.seq))
	binary.PutUvarint(buf[12:16], uint64(msg.ack))
	binary.PutUvarint(buf[16:20], uint64(msg.channel))
}

func ParseMessageHeader(buf []byte) (*Message, uint16) {
	msg := &Message{
		syn:     buf[6] > 0,
		fin:     buf[7] > 0,
		seq:     bytesToUint16(buf[8:12]),
		ack:     bytesToUint16(buf[12:16]),
		channel: bytesToUint16(buf[16:20]),
	}

	return msg, bytesToUint16(buf[:4])
}
