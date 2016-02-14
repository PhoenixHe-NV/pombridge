package main

import (
	"net"
	bridge "pombridge"
	"pombridge/core"
	"pombridge/log"
)

var server = bridge.NewServer()

func main() {
	run()
}

func run() {
	err := server.Listen(bridge.Addr{"127.0.0.1", 8000})
	if err != nil {
		log.F("Bridge listen: ", err)
	}

	log.I("Start listening")

	for {
		conn, err := server.Accept()
		if err != nil {
			log.E("Bridge accept: ", err)
			continue
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	log.I("Bridge accept: ", conn.RemoteAddr())

	remote, err := net.Dial("tcp", "127.0.0.1:1080")
	if err != nil {
		log.W("Cannot dail remote ", err)
		conn.Close()
		return
	}

	go core.ConnCopy(conn, remote)
	core.ConnCopy(remote, conn)

	log.I("Bridge closed: ", conn.RemoteAddr())
}
