package mac

import (
	"encoding/hex"
	"fmt"
	"net"
)

const (
	MINMACLENGTH = 12
	MAXMACLENGTH = 59
	MAXBUFCAP    = 48

	MACLENGTH12 = 12
	MACLENGTH16 = 16
	MACLENGTH40 = 40
)

func ParseMAC(mac string) (net.HardwareAddr, error) {
	if l := len(mac); l < MINMACLENGTH || l > MAXMACLENGTH {
		return nil, fmt.Errorf("mac %s is invalid with length %d not in [12, 59]", mac, l)
	}

	buf := make([]byte, 0, MAXBUFCAP)
	count := 0
	for _, r := range mac {
		switch {
		case r >= '0' && r <= '9':
			buf = append(buf, byte(r))
			count++
		case r >= 'A' && r <= 'F':
			buf = append(buf, byte(r))
			count++
		case r >= 'a' && r <= 'f':
			buf = append(buf, byte(r))
			count++
		default:
			continue
		}

		if count >= MAXBUFCAP {
			break
		}
	}

	switch count {
	case MACLENGTH12, MACLENGTH16, MACLENGTH40:
	default:
		return nil, fmt.Errorf("mac %s is invalid with chars length %d not in [12, 16, 40]", mac, count)
	}

	hwaddr := make([]byte, count/2)
	if _, err := hex.Decode(hwaddr, buf[:count]); err != nil {
		return nil, fmt.Errorf("parse mac %s failed: %s", mac, err.Error())
	}

	return net.HardwareAddr(hwaddr), nil
}

func Mac48ToUint64(mac net.HardwareAddr) uint64 {
	if len(mac) == 0 {
		return 0
	} else {
		return (uint64(mac[0]) | uint64(mac[1])<<8 | uint64(mac[2])<<16 | uint64(mac[3])<<24 |
			uint64(mac[4])<<32 | uint64(mac[5])<<40)
	}
}

func Mac48FromEUI64(ip net.IP) (net.HardwareAddr, error) {
	if ip.To16() == nil {
		return nil, fmt.Errorf("ip address shorter than 16 bytes")
	}

	if isEUI48 := ip[11] == 0xff && ip[12] == 0xfe; !isEUI48 {
		return nil, fmt.Errorf("ip address is not an EUI48 address")
	}

	mac := make(net.HardwareAddr, 6)
	copy(mac[0:3], ip[8:11])
	copy(mac[3:6], ip[13:16])
	mac[0] ^= 0x02
	return mac, nil
}

func Mac48ToEUI64(ip net.IP, mac net.HardwareAddr) (net.IP, error) {
	if ip.To4() != nil {
		return nil, fmt.Errorf("ip address must not be ipv4")
	}

	if len(mac) != 6 {
		return nil, fmt.Errorf("mac address must be 6 bytes")
	}

	ip16 := ip.To16()
	if ip16 == nil {
		return nil, fmt.Errorf("ip address shorter than 16 bytes")
	}

	ipv6 := make(net.IP, 16)
	copy(ipv6, ip16)
	copy(ipv6[8:11], mac[0:3])
	copy(ipv6[13:16], mac[3:6])
	ipv6[8] ^= 0x02
	ipv6[11] = 0xff
	ipv6[12] = 0xfe
	return ipv6, nil
}
