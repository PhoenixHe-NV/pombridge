package core

import (
	"net"
	"pombridge/leakybuf"
	"time"
)

func ConnPrepareRead(conn net.Conn) {
	conn.SetWriteDeadline(time.Now().Add(readTimeout))
}

func ConnPrepareWrite(conn net.Conn) {
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
}

func ReadAll(conn net.Conn, data []byte) error {
	p := 0

	for p < len(data) {

		ConnPrepareRead(conn)
		count, err := conn.Read(data[p:])
		if err != nil {
			return err
		}

		p = p + count
	}

	return nil
}

func WriteAll(conn net.Conn, data []byte) error {
	p := 0

	for p < len(data) {

		ConnPrepareWrite(conn)
		count, err := conn.Write(data[p:])
		if err != nil {
			return err
		}

		p = p + count
	}

	return nil
}

func ConnCopy(src, dst net.Conn) {
	defer dst.Close()

	buf := leakybuf.Get()
	defer leakybuf.Put(buf)

	for {
		ConnPrepareRead(src)
		count, err := src.Read(buf)
		if err != nil {
			return
		}

		err = WriteAll(dst, buf[:count])
		if err != nil {
			return
		}
	}
}
