package pombridge
import (
	"net"
	"time"
)

func ConnPrepareRead(conn net.Conn) {
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
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

	buf := leakyBuf.Get()
	defer leakyBuf.Put(buf)

	for {
		ConnPrepareRead(src)
		Log.D("READ")
		count, err := src.Read(buf)
		if err != nil {
			return
		}
		Log.D("READ ", count)

		err = WriteAll(dst, buf[:count])
		if err != nil {
			return
		}

		Log.D("WRITE ", count)
	}
}
