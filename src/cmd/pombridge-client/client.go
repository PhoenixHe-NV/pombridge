package main

import (
	"flag"
	"fmt"
	"net"
	bridge "pombridge"
	"pombridge/core"
	"pombridge/log"
)

var client = bridge.NewClient()

var config struct {
	listenPort int
}

func flagInit() {
	flag.IntVar(&config.listenPort, "p", 0, "client listen port")

	flag.Parse()
}

func main() {
	flagInit()

	client.Connect(bridge.Addr{"127.0.0.1", 8000})
	run("127.0.0.1:" + fmt.Sprint(config.listenPort))
}

func run(listenAddr string) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.F(err)
	}

	log.I("Start listening at " + listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.E("accept: ", err)
			continue
		}

		log.I("Accept connection " + conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	remote, err := client.Dial()
	if err != nil {
		log.W("Cannot dail remote")
		conn.Close()
		return
	}

	go core.ConnCopy(conn, remote)
	core.ConnCopy(remote, conn)

	log.I("Close connection " + conn.RemoteAddr().String())
}
