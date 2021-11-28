package binary

import (
	encodingbinary "encoding/binary"
	"fmt"
	"io"
	"net"
)

func Exchange(conn net.Conn, req []byte, ignoreResp bool) ([]byte, error) {
	if err := Write(conn, req); err != nil {
		return nil, err
	} else if ignoreResp {
		return nil, nil
	} else {
		return Read(conn)
	}
}

func Write(conn net.Conn, req []byte) error {
	if err := encodingbinary.Write(conn, encodingbinary.BigEndian, uint16(len(req))); err != nil {
		return fmt.Errorf("write msg size failed: %s", err.Error())
	} else if _, err = conn.Write(req); err != nil {
		return fmt.Errorf("write msg body failed: %s", err.Error())
	} else {
		return nil
	}
}

func Read(conn net.Conn) ([]byte, error) {
	var size uint16
	if err := encodingbinary.Read(conn, encodingbinary.BigEndian, &size); err != nil {
		return nil, fmt.Errorf("read msg size failed: %s", err.Error())
	} else {
		buf := make([]byte, size)
		if _, err = io.ReadFull(conn, buf); err != nil {
			return nil, fmt.Errorf("read msg body failed: %s", err.Error())
		} else {
			return buf, nil
		}
	}
}
