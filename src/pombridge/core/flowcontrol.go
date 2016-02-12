package core

import (
	"sync"
	"pombridge/heap"
	"math"
)

type Bridge struct {

}

type BridgeFlow struct {
	mutex *sync.RWMutex
	seq, ack uint16
}

var FlowControl = &BridgeFlow{
	mutex: &sync.RWMutex{},
	seq: 0,
	ack: 0,
}

var recv = make(MsgChan)

func (f *BridgeFlow) NextSeq() uint16 {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.seq = f.seq + 1
	return f.seq
}

func (f *BridgeFlow) Ack() uint16 {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.ack
}

func (f *BridgeFlow) setAck(ack uint16) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.ack = ack
}

func (f *BridgeFlow) RecvMsg(msg *Message) {
	recv <- msg
}

func (msg *Message) Priority() int {
	// Do not need to lock because only the goroutine that
	// runs runFlowControl() will modify heap and read
	// the priority

	p := int(msg.seq)
	// This algorithm works because we only need to guarantee
	// the comparisons between priority are correct
	if p < int(FlowControl.ack) {
		p = p + math.MaxInt16
	}

	return p
}

func runFlowControl() {
	heap := heap.New()

	for {
		msg := <- recv
		heap.Push(msg)

		ack := FlowControl.ack
		for !heap.Empty() {
			msg := heap.Top().(*Message)
			if msg.seq != ack + 1 {
				break
			}

			heap.Pop()
			ack = msg.seq
			ch, ok := BusChannel(msg.channel)
			if !ok {
				// channel not found, ignored
				break
			}

			ch <- msg
		}
	}
}
