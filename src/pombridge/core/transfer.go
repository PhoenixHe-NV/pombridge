package core

import (
	"net"
	"pombridge/leakybuf"
	"pombridge/log"
	"sync"
)

type Bridge struct {
	busesMutex      *sync.RWMutex
	flow            *Flow
	SendBus         MsgChan
	highPrioSendBus MsgChan
	recvBuses       map[uint16](MsgChan) // map channel to wait chan
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
	msg := b.flow.NewMsgToSend(channel, nil)
	msg.syn = true
	b.SendBus <- msg
}

func (b *Bridge) CloseChannel(channel uint16) {
	b.busesMutex.Lock()
	defer b.busesMutex.Unlock()

	delete(b.recvBuses, channel)
	msg := b.flow.NewMsgToSend(channel, nil)
	msg.fin = true
	b.SendBus <- msg
}

func (b *Bridge) BusChannel(channel uint16) (MsgChan, bool) {
	b.busesMutex.RLock()
	defer b.busesMutex.RUnlock()

	ch, ok := b.recvBuses[channel]
	return ch, ok
}

func (bridge *Bridge) RunBusLine(conn net.Conn) {
	connClosed := make(chan int)
	go bridge.RunBusLineRecv(conn, connClosed)
	bridge.RunBusLineSend(conn, connClosed)
	close(connClosed)
}

func (bridge *Bridge) RunBusLineSend(conn net.Conn, connClosed chan int) {
	buf := leakybuf.Get()
	defer leakybuf.Put(buf)
	var msg *Message = nil

LOOP:
	for {
		select {
		case msg = <-bridge.highPrioSendBus:
		case msg = <-bridge.SendBus:
		case <-connClosed:
			break LOOP
		}

		msg.ack = msg.flow.Ack()
		msg.PacketHeader(buf)
		log.D("msg > seq:", msg.seq, " ack:", msg.ack,
			" channel:", msg.channel, " syn:", msg.syn, " fin:", msg.fin)
		log.D("msg > ", string(msg.data))
		err := WriteAll(conn, buf[:MsgeaderLen])
		if err != nil {
			log.D("BuslineSend: ", err)
			break
		}
		err = WriteAll(conn, msg.data)
		if err != nil {
			log.D("BuslineSend: ", err)
			break
		}

		msg = nil
	}

	conn.Close()
	if msg != nil {
		// resend the message
		bridge.highPrioSendBus <- msg
	}
}

func (bridge *Bridge) RunBusLineRecv(conn net.Conn, connClosed chan int) {
	buf := leakybuf.Get()
	defer leakybuf.Put(buf)

	for {
		log.D("ready to recv header in busline")
		err := ReadAll(conn, buf[:MsgeaderLen])
		if err != nil {
			log.D("BuslineRecv: ", err)
			break
		}
		msg, dataLen := ParseMessageHeader(buf)
		msg.flow = bridge.flow
		msg.data = make([]byte, dataLen)
		log.D("ready to recv body in busline, len:", dataLen)
		err = ReadAll(conn, msg.data)
		if err != nil {
			log.D("BuslineRecv: ", err)
			break
		}

		msg.flow.RecvMsg(msg)
	}

	connClosed <- 1
	conn.Close()
}
