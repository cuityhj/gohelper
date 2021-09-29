package binary

import (
	encodingbinary "encoding/binary"
	"fmt"
	"io"
	"net"
)

func Exchange(conn net.Conn, req []byte) ([]byte, error) {
	if err := Write(conn, req); err != nil {
		return nil, fmt.Errorf("write msg to server failed: %s", err.Error())
	} else {
		return Read(conn)
	}
}

func Write(conn net.Conn, req []byte) error {
	err := encodingbinary.Write(conn, encodingbinary.BigEndian, uint16(len(req)))
	if err != nil {
		return fmt.Errorf("write msg size to server failed: %s", err.Error())
	}

	_, err = conn.Write(req)
	return err
}

func Read(conn net.Conn) ([]byte, error) {
	var size uint16
	err := encodingbinary.Read(conn, encodingbinary.BigEndian, &size)
	if err != nil {
		return nil, fmt.Errorf("read msg size from server failed: %s", err.Error())
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(conn, buf)
	return buf, err
}
