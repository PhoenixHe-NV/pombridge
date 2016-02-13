package core

import (
	"net"
	"pombridge/leakybuf"
	"sync"
)

type Bridge struct {
	busesMutex      *sync.RWMutex
	flow            *Flow
	SendBus         MsgChan
	highPrioSendBus MsgChan
	recvBuses       map[uint16](MsgChan) // map channel to wait chan
}

type BusLine struct {
	bridge   *Bridge
	conn     net.Conn
	recvChan chan Message
}

func NewBridge() *Bridge {
	b := &Bridge{
		busesMutex:      &sync.RWMutex{},
		SendBus:         make(MsgChan),
		highPrioSendBus: make(MsgChan),
		recvBuses:       make(map[uint16](MsgChan)),
	}
	b.initFlow()

	return b
}

func (b *Bridge) OpenChannel(channel uint16, recvChan MsgChan) {
	b.busesMutex.Lock()
	defer b.busesMutex.Unlock()

	b.recvBuses[channel] = recvChan
}

func (b *Bridge) CloseChannel(channel uint16) {
	b.busesMutex.Lock()
	defer b.busesMutex.Unlock()

	delete(b.recvBuses, channel)
}

func (b *Bridge) BusChannel(channel uint16) (MsgChan, bool) {
	b.busesMutex.RLock()
	defer b.busesMutex.RUnlock()

	ch, ok := b.recvBuses[channel]
	return ch, ok
}

func (b *BusLine) runSend() {
	buf := leakybuf.Get()
	defer leakybuf.Put(buf)
	var msg *Message = nil

	for {
		select {
		case msg = <-b.bridge.highPrioSendBus:
		case msg = <-b.bridge.SendBus:
		}

		msg.ack = msg.flow.Ack()
		msg.PacketHeader(buf)
		err := WriteAll(b.conn, buf[:PacketHeaderLen])
		if err != nil {
			break
		}
		err = WriteAll(b.conn, msg.data)
		if err != nil {
			break
		}

		msg = nil
	}

	b.conn.Close()
	if msg != nil {
		// resend the message
		b.bridge.highPrioSendBus <- msg
	}
}

func (b *BusLine) runRecv() {
	buf := leakybuf.Get()
	defer leakybuf.Put(buf)

	for {
		err := ReadAll(b.conn, buf[:PacketHeaderLen])
		if err != nil {
			break
		}
		msg, dataLen := ParseMessageHeader(buf)
		msg.data = make([]byte, dataLen)
		err = ReadAll(b.conn, msg.data)
		if err != nil {
			break
		}

		msg.flow.RecvMsg(msg)
	}

	b.conn.Close()
}
