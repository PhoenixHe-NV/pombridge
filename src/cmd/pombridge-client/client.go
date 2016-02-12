package main

import (
	pb "pombridge"
	"flag"
	"net"
	"fmt"
)

var Log *pb.PomLogger = &pb.Log

var config struct {
	listenPort int
}

func flagInit() {
	flag.IntVar(&config.listenPort, "p", 0, "client listen port")

	flag.Parse()
}

func main() {
	pb.CommonInit()
	flagInit()

	run("127.0.0.1:" + fmt.Sprint(config.listenPort))
}

func run(listenAddr string) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		Log.F(err)
	}

	Log.I("Start listening at " + listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			Log.E("accept: ", err)
			continue
		}

		Log.I("Accept connection " + conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	remote, err := pb.Dial()
	if err != nil {
		Log.W("Cannot dail remote")
		conn.Close()
		return
	}

	go pb.ConnCopy(conn, remote)
	pb.ConnCopy(remote, conn)

	Log.I("Close connection " + conn.RemoteAddr().String())
}
