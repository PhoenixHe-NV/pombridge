package pombridge

import (
	"fmt"
	"net"
	"pombridge/core"
	"pombridge/log"
	"time"
)

type Addr struct {
	Host string
	Port uint16
}

// TODO use cookie to identify different clients

type Client struct {
	bridge *core.Bridge
}

type Server struct {
	bridge *core.Bridge
}

func (addr Addr) String() string {
	return addr.Host + ":" + fmt.Sprint(addr.Port)
}

func NewClient() *Client {
	return &Client{core.NewBridge()}
}

func NewServer() *Server {
	return &Server{core.NewBridge()}
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
	return core.NewChannel(c.bridge), nil
}

func (s *Server) Listen(addr Addr) error {
	listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

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
	time.Sleep(time.Hour)
	return nil, nil
}
