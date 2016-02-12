package core

import (
	"net"
	"time"
	"errors"
)

type Channel struct {
	bridge Bridge
	id uint16
	ch chan []byte
	closed bool
}

func NewChannel() *Channel {
	return &Channel{ch: make(chan []byte)}
}

func (c* Channel) Read(b []byte) (int, error) {
	if (c.closed) {
		return 0, errors.New("bridge closed")
	}

	buf := <- c.ch
	copy(b, buf)
	return len(buf), nil
}

func (c* Channel) Write(b []byte) (int, error) {
	if (c.closed) {
		return 0, errors.New("bridge closed")
	}

	buf := make([]byte, len(b))
	copy(buf, b)
	c.ch <- buf
	return len(b), nil
}

func (c* Channel) Close() error {
	close(c.ch)
	c.closed = true
	return nil
}

func (c* Channel) LocalAddr() net.Addr {
	return nil
}

func (c* Channel) RemoteAddr() net.Addr {
	return nil
}

func (c* Channel) SetDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c* Channel) SetReadDeadline(t time.Time) error {
	return nil
}

func (c* Channel) SetWriteDeadline(t time.Time) error {
	return nil
}
