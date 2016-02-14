package core

import (
	"math"
	"pombridge/heap"
	"pombridge/log"
	"sync"
)

type Flow struct {
	bridge   *Bridge
	mutex    *sync.RWMutex
	seq, ack uint16
	recv     MsgChan
}

func (b *Bridge) initFlow() {
	b.flow = &Flow{
		bridge: b,
		mutex:  &sync.RWMutex{},
		seq:    0,
		ack:    0,
		recv:   make(MsgChan),
	}
	go b.flow.runFlowControl()
}

func (f *Flow) NewMsgToSend(channel uint16, b []byte) *Message {
	msg := &Message{
		flow:    f,
		syn:     false,
		fin:     false,
		seq:     f.NextSeq(),
		channel: channel,
		data:    make([]byte, len(b)),
	}
	copy(msg.data, b)

	return msg
}

func (f *Flow) NextSeq() uint16 {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.seq = f.seq + 1
	return f.seq
}

func (f *Flow) Ack() uint16 {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.ack
}

func (f *Flow) setAck(ack uint16) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.ack = ack
}

func (f *Flow) RecvMsg(msg *Message) {
	log.D("msg < seq:", msg.seq, " ack:", msg.ack,
		" channel:", msg.channel, " syn:", msg.syn, " fin:", msg.fin)
	log.D("msg < ", string(msg.data))
	f.recv <- msg
}

func (msg *Message) Priority() int {
	// Do not need to lock because only the goroutine that
	// runs runFlowControl() will modify heap and read
	// the priority

	p := int(msg.seq)
	// This algorithm works because we only need to guarantee
	// the comparisons between priority are correct
	if p < int(msg.flow.ack) {
		p = p + math.MaxInt16
	}

	return p
}

// TODO add chan to inform bridge closed
func (f *Flow) runFlowControl() {
	heap := heap.New()

	for {
		msg := <-f.recv
		heap.Push(msg)

		ack := f.ack
		for !heap.Empty() {
			msg := heap.Top().(*Message)
			if msg.seq != ack+1 {
				break
			}

			heap.Pop()
			ack = msg.seq

			ch, ok := f.bridge.BusChannel(msg.channel)
			if !ok {
				if msg.syn {
					// new channel to accept
					if f.bridge.canAccept {
						conn := f.bridge.NewChannel(msg.channel)
						select {
						case f.bridge.AcceptChan <- conn:
							log.I("new channel opened: ", msg.channel)
						default:
							conn.Close()
							log.I("new channel denied: ", msg.channel)
						}
						continue
					}
				}
				// channel not found, ignored
				log.I("ignore packet from unknown channel: ", msg.channel)
				continue
			}

			log.D("dispatch msg to channel ", msg.channel)
			trySend(ch, msg)
		}

		f.setAck(ack)
		// TODO set timeout to send extra ack packet
	}
}

func trySend(ch MsgChan, msg *Message) {
	defer func() { recover() }()
	ch <- msg
}
