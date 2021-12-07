package mac

import "net"

func Mac48ToUint64(mac net.HardwareAddr) uint64 {
	return (uint64(mac[0]) | uint64(mac[1])<<8 | uint64(mac[2])<<16 | uint64(mac[3])<<24 |
		uint64(mac[4])<<32 | uint64(mac[5])<<40)
}
