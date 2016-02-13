package main

import (
	"net"
	bridge "pombridge"
	"pombridge/core"
	"pombridge/log"
)

var server = bridge.NewServer()

func main() {
}

func run() {
	err := server.Listen(bridge.Addr{"127.0.0.1", 8000})
	if err != nil {
		log.F("bridge listen: ", err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			log.E("bridge accept: ", err)
			continue
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	remote, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.W("Cannot dail remote")
		conn.Close()
		return
	}

	go core.ConnCopy(conn, remote)
	core.ConnCopy(remote, conn)
}
