package pombridge
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
	syn, fin bool
	seq, ack uint16
	channel uint16
	data []byte
}

type MsgChan chan *Message

const (
	versionMajor = 1
	versionMinor = 0
	PacketHeaderLen = 20
)

func (msg *Message) PacketHeader(buf []byte) {
	packetLen := PacketHeaderLen + len(msg.data)

	binary.PutUvarint(buf[0:4], packetLen)
	buf[4] = versionMajor
	buf[5] = versionMinor
	buf[6] = msg.syn
	buf[7] = msg.fin
	binary.PutUvarint(buf[8:12], msg.seq)
	binary.PutUvarint(buf[12:16], msg.ack)
	binary.PutUvarint(buf[16:20], msg.channel)

	return buf
}

func ParseMessageHeader(buf []byte) (*Message, uint16) {
	msg := Message{
		syn: bool(buf[6]),
		ack: bool(buf[7]),
		seq: binary.Uvarint(buf[8:12]),
		ack: binary.Uvarint(buf[12:16]),
		channel: binary.Uvarint(buf[16:20]),
	}

	return &msg, binary.Uvarint(buf[:4])
}
