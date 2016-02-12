package pombridge
import (
	"net"
	"time"
	"errors"
)

type BridgeAddr struct {
	host string
	port int16
}

type Bridge struct {
	id int32
	seg, ack int32
	sendBuf []byte
	ch chan []byte
	closed bool
}

func (bridge* Bridge) Read(b []byte) (int, error) {
	if (bridge.closed) {
		return 0, errors.New("bridge closed")
	}

	buf := <- bridge.ch
	copy(b, buf)
	return len(buf), nil
}

func (bridge* Bridge) Write(b []byte) (int, error) {
	if (bridge.closed) {
		return 0, errors.New("bridge closed")
	}

	buf := make([]byte, len(b))
	copy(buf, b)
	bridge.ch <- buf
	return len(b), nil
}

func (bridge* Bridge) Close() error {
	close(bridge.ch)
	bridge.closed = true
	return nil
}

func (bridge* Bridge) LocalAddr() net.Addr {
	return nil
}

func (bridge* Bridge) RemoteAddr() net.Addr {
	return nil
}

func (bridge* Bridge) SetDeadline(t time.Time) error {
	err := bridge.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return bridge.SetWriteDeadline(t)
}

func (bridge* Bridge) SetReadDeadline(t time.Time) error {
	return nil
}

func (bridge* Bridge) SetWriteDeadline(t time.Time) error {
	return nil
}

func SetDialAddr(remotes []BridgeAddr) {

}

func Dial() (net.Conn, error) {
	return &Bridge{
		ch: make(chan []byte),
	}, nil
}

func Listen() error {
	return nil
}

func Accept() (net.Conn, error) {
	return nil, nil
}