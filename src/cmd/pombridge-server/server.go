package main

import (
	pb "pombridge"
	"net"
)

var Log *pb.PomLogger = &pb.Log

func main() {
	pb.CommonInit()

}

func run() {
	err := pb.Listen()
	if err != nil {
		Log.F("bridge listen: ", err)
	}

	for {
		conn, err := pb.Accept()
		if (err != nil) {
			Log.E("bridge accept: ", err)
			continue
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	remote, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		Log.W("Cannot dail remote")
		conn.Close()
		return
	}

	go pb.ConnCopy(conn, remote)
	pb.ConnCopy(remote, conn)
}
