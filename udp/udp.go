package udp

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/cuityhj/gohelper/syncpool"
)

var DefaultTimeout = 10 * time.Second

var gDefaultUDPClient = &UDPClient{timeout: DefaultTimeout}

func GetDefaultUDPClient() *UDPClient {
	return gDefaultUDPClient
}

type UDPClient struct {
	timeout time.Duration
}

func NewUDPClient(timeout time.Duration) *UDPClient {
	return &UDPClient{timeout: timeout}
}

func (cli *UDPClient) Exchange(request []byte, serverIp net.IP, serverPort uint32, response *bytes.Buffer) error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: serverIp, Port: int(serverPort)})
	if err != nil {
		return fmt.Errorf("dial udp with addr %s and port %d failed: %s", serverIp, serverPort, err.Error())
	}

	defer conn.Close()

	if _, err := conn.Write(request); err != nil {
		return fmt.Errorf("write udp request to addr %s and port %d failed: %s", serverIp, serverPort, err.Error())
	}

	conn.SetReadDeadline(time.Now().Add(cli.timeout))
	buf := syncpool.GetBytePool().Get()
	defer syncpool.GetBytePool().Put(buf)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("read udp respsonse from addr %s and port %d failed: %s", serverIp, serverPort, err.Error())
	}

	if response != nil {
		response.Write(buf[:n])
	}

	return nil
}

func (cli *UDPClient) Write(request []byte, serverIp net.IP, serverPort uint32) error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: serverIp, Port: int(serverPort)})
	if err != nil {
		return fmt.Errorf("dial udp with addr %s and port %d failed: %s", serverIp, serverPort, err.Error())
	}

	defer conn.Close()
	if _, err := conn.Write(request); err != nil {
		return fmt.Errorf("write udp request to addr %s and port %d failed: %s", serverIp, serverPort, err.Error())
	}

	return nil
}
