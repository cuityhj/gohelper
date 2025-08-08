package ip

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"net/netip"
)

type IP net.IP

func (one IP) Cmp(another IP) int {
	oneAddr, _ := netip.AddrFromSlice(one)
	anotherAddr, _ := netip.AddrFromSlice(another)
	return oneAddr.Compare(anotherAddr)
}

func IsIpZero(ip net.IP) bool {
	return ip == nil || ip.IsUnspecified()
}

func CheckIPsValid(ips ...string) error {
	for _, ip := range ips {
		if net.ParseIP(ip) == nil {
			return fmt.Errorf("invalid ip %s", ip)
		}
	}

	return nil
}

func CheckIPv4sValid(ips ...string) error {
	for _, ip := range ips {
		if _, err := ParseIP(ip, true); err != nil {
			return err
		}
	}

	return nil
}

func CheckIPv6sValid(ips ...string) error {
	for _, ip := range ips {
		if _, err := ParseIP(ip, false); err != nil {
			return err
		}
	}

	return nil
}

func ParseIPv4(ipstr string) (net.IP, error) {
	return ParseIP(ipstr, true)
}

func ParseIPv6(ipstr string) (net.IP, error) {
	return ParseIP(ipstr, false)
}

func ParseIP(ipstr string, isv4 bool) (net.IP, error) {
	if ip := net.ParseIP(ipstr); ip == nil || (ip.To4() != nil) != isv4 {
		return nil, fmt.Errorf("invalid ip %s", ipstr)
	} else {
		if isv4 {
			return ip.To4(), nil
		} else {
			return ip, nil
		}
	}
}

func IPv4FromUint32(i uint32) net.IP {
	ip := make([]byte, 4)
	binary.BigEndian.PutUint32(ip, i)
	return net.IP(ip).To4()
}

func IPv4ToUint32(ip net.IP) uint32 {
	if ip == nil {
		return 0
	} else {
		return binary.BigEndian.Uint32(ip.To4())
	}
}

func IPv4StrToUint32(ipstr string) (uint32, error) {
	if ipv4, err := ParseIP(ipstr, true); err != nil {
		return 0, err
	} else {
		return IPv4ToUint32(ipv4), nil
	}
}

func IPv6FromBigInt(bigint *big.Int) net.IP {
	var ip net.IP
	if bigint == nil || len(bigint.Bytes()) > net.IPv6len {
		return ip
	}

	bytes := bigint.Bytes()
	bytesLen := len(bigint.Bytes())
	for i := 0; i < net.IPv6len-bytesLen; i++ {
		bytes = append([]byte{0}, bytes...)
	}

	return net.IP(bytes)
}

func IPv6ToBigInt(ip net.IP) *big.Int {
	return new(big.Int).SetBytes(ip.To16())
}

func IPv6StrToBigInt(ipstr string) (*big.Int, error) {
	if ipv6, err := ParseIP(ipstr, false); err != nil {
		return nil, err
	} else {
		return IPv6ToBigInt(ipv6), nil
	}
}

func ParseIPAndVersion(ipstr string) (net.IP, bool, error) {
	if ip := net.ParseIP(ipstr); ip == nil {
		return nil, false, fmt.Errorf("invalid ip %s", ipstr)
	} else if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4, false, nil
	} else {
		return ip, true, nil
	}
}
