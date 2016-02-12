package core

import (
	"net"
	"sync"
	"pombridge/leakybuf"
)

var SendBus = make(MsgChan)
var highPrioSendBus = make(MsgChan)
var recvBuses = make(map[uint16](MsgChan))	// map channel to wait chan
var busesMutex = &sync.RWMutex{}


type BusLine struct {
	conn net.Conn
	recvChan chan Message
}

func OpenChannel(channel uint16, recvChan MsgChan) {
	busesMutex.Lock()
	defer busesMutex.Unlock()

	recvBuses[channel] = recvChan
}

func CloseChannel(channel uint16) {
	busesMutex.Lock()
	defer busesMutex.Unlock()

	delete(recvBuses, channel)
}

func BusChannel(channel uint16) (MsgChan, bool) {
	busesMutex.RLock()
	defer busesMutex.RUnlock()

	ch, ok := recvBuses[channel]
	return ch, ok
}

func (busLine *BusLine) runSend() {
	buf := leakybuf.Get()
	defer leakybuf.Put(buf)
	var msg *Message = nil

	for {
		select {
		case msg = <- highPrioSendBus:
		case msg = <- SendBus:
		}

		msg.seq = FlowControl.NextSeq()
		msg.ack = FlowControl.Ack()
		msg.PacketHeader(buf)
		err := WriteAll(busLine.conn, buf[:PacketHeaderLen])
		if err != nil {
			break
		}
		err = WriteAll(busLine.conn, msg.data)
		if err != nil {
			break
		}

		msg = nil
	}

	busLine.conn.Close()
	if msg != nil {
		// resend the message
		highPrioSendBus <- msg
	}
}

func (busLine *BusLine) runRecv() {
	buf := leakybuf.Get()
	defer leakybuf.Put(buf)

	for {
		err := ReadAll(busLine.conn, buf[:PacketHeaderLen])
		if err != nil {
			break
		}
		msg, dataLen := ParseMessageHeader(buf)
		msg.data = make([]byte, dataLen)
		err = ReadAll(busLine.conn, msg.data)
		if err != nil {
			break
		}

		FlowControl.RecvMsg(msg)
	}

	busLine.conn.Close()
}



