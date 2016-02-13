package pombridge

import (
	"net"
	"pombridge/core"
)

type Addr struct {
	Host string
	Port uint16
}

type Client struct {
	core.Bridge
}

type Server struct {
	core.Bridge
}

func NewClient() *Client {
	return &Client{}
}

func NewServer() *Server {
	return &Server{}
}

func (c *Client) Connect(addr Addr) error {
	return nil
}

func (c *Client) Dial() (net.Conn, error) {
	return core.NewChannel(), nil
}

func (s *Server) Listen(addr Addr) error {
	return nil
}

func (s *Server) Accept() (net.Conn, error) {
	return nil, nil
}
