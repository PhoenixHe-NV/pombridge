package core

import (
	"errors"
	"io"
	"net"
	"pombridge/log"
	"sync"
	"time"
)

type Channel struct {
	mutex  *sync.RWMutex
	bridge *Bridge
	id     uint16
	closed bool
	recv   chan *Message
}

func NewChannel(bridge *Bridge) *Channel {
	c := &Channel{
		mutex:  &sync.RWMutex{},
		bridge: bridge,
		id:     0,
		recv:   make(chan *Message),
	}
	bridge.OpenChannel(0, c.recv)
	return c
}

func (c *Channel) Read(buf []byte) (int, error) {
	if c.Closed() {
		return 0, io.ErrClosedPipe
	}

	msg := <-c.recv
	if msg == nil || msg.fin {
		c.Close()
		return 0, io.ErrClosedPipe
	}
	if len(msg.data) > len(buf) {
		log.E("Recevive message which size is bigger than excepted!")
		return 0, errors.New("unexcepted message")
	}
	copy(buf, msg.data)

	return len(msg.data), nil
}

func (c *Channel) Write(b []byte) (int, error) {
	if c.Closed() {
		return 0, io.ErrClosedPipe
	}

	msg := c.bridge.flow.NewMsgToSend(c.id, b)
	c.bridge.SendBus <- msg

	return len(b), nil
}

func (c *Channel) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return io.ErrClosedPipe
	}

	close(c.recv)
	c.bridge.CloseChannel(c.id)
	c.closed = true
	return nil
}

func (c *Channel) Closed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.closed
}

func (c *Channel) LocalAddr() net.Addr {
	return nil
}

func (c *Channel) RemoteAddr() net.Addr {
	return nil
}

func (c *Channel) SetDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c *Channel) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Channel) SetWriteDeadline(t time.Time) error {
	return nil
}
