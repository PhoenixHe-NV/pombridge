package core

import (
	"io"
	"net"
	"pombridge/log"
	"time"
)

type Channel struct {
	bridge *Bridge
	id     uint16
	closed bool
	recv   MsgChan
}

func NewChannel(bridge *Bridge) *Channel {
	c := &Channel{
		bridge: bridge,
		id:     0,
		recv:   make(MsgChan),
	}
	bridge.OpenChannel(0, c.recv)
	return c
}

func (c *Channel) Read(buf []byte) (int, error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}

	msg := <-c.recv
	if len(msg.data) > len(buf) {
		log.F("Recevive a packet which size is bigger than excepted!")
	}

	return len(buf), nil
}

func (c *Channel) Write(b []byte) (int, error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}

	msg := c.bridge.flow.NewMsgToSend(c.id, b)
	c.bridge.SendBus <- msg

	return len(b), nil
}

func (c *Channel) Close() error {
	if c.closed {
		return io.ErrClosedPipe
	}

	close(c.recv)
	c.bridge.CloseChannel(c.id)
	c.closed = true
	return nil
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
