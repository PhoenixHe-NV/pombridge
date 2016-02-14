package pombridge

import (
	"errors"
	"fmt"
	"math"
	"net"
	"pombridge/core"
	"pombridge/log"
	"sync/atomic"
	"time"
)

type Addr struct {
	Host string
	Port uint16
}

// TODO use cookie to identify different clients
type Client struct {
	bridge        *core.Bridge
	freeChannelId chan uint16
	maxChannelId  uint32
}

type Server struct {
	bridge    *core.Bridge
	listening bool
}

func (addr Addr) String() string {
	return addr.Host + ":" + fmt.Sprint(addr.Port)
}

func NewClient() *Client {
	return &Client{
		bridge:        core.NewBridge(0),
		freeChannelId: make(chan uint16, math.MaxUint16),
	}
}

func NewServer() *Server {
	return &Server{
		bridge:    core.NewBridge(5),
		listening: false,
	}
}

func (c *Client) Connect(addr Addr) {
	go c.runBusLine(addr)
}

func (c *Client) runBusLine(addr Addr) {
	// TODO use config param
	// TODO use 2 or more connections to one addr
	for {
		conn, err := net.Dial("tcp", addr.String())
		if err != nil {
			log.E("Busline failed: ", err)
			time.Sleep(time.Second * 2)
			continue
		}

		log.I("Busline established: ", addr)

		c.bridge.RunBusLine(conn)

		log.I("Busline closed: ", addr, " , reconnect")
	}
}

func (c *Client) Dial() (net.Conn, error) {
	id := c.genNewChannelId()
	conn := c.bridge.NewChannel(id)
	conn.AfterClose(func(ch *core.Channel) {
		c.freeChannelId <- id
	})
	return conn, nil
}

func (c *Client) genNewChannelId() uint16 {
	select {
	case id := <-c.freeChannelId:
		return id
	default:
		id := atomic.AddUint32(&c.maxChannelId, 1)
		if id > math.MaxUint16 {
			log.E("Too many channels! ", id, " > ", math.MaxUint16)
		}
		return uint16(id)
	}
}

func (s *Server) Listen(addr Addr) error {
	listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

	s.listening = true

	go s.runBusLine(listener)
	return nil
}

func (s *Server) runBusLine(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.E("Busline accept failed: ", err)
			continue
		}

		go s.serveBusLine(conn)
	}
}

func (s *Server) serveBusLine(conn net.Conn) {
	log.I("Busline opened: ", conn.RemoteAddr())

	s.bridge.RunBusLine(conn)

	log.I("Busline closed: ", conn.RemoteAddr())
}

func (s *Server) Accept() (net.Conn, error) {
	if !s.listening {
		return nil, errors.New("bridge server is not listening")
	}
	conn := <-s.bridge.AcceptChan
	return conn, nil
}
