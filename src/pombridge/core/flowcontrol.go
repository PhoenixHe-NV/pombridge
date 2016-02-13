package core

import (
	"math"
	"pombridge/heap"
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
				// channel not found, ignored
				break
			}

			ch <- msg
		}
	}
}
