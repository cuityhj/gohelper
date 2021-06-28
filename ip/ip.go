package ip

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
)

func IsIpZero(ip net.IP) bool {
	return ip == nil || ip.IsUnspecified()
}

func ParseIPv4(ipstr string) (net.IP, error) {
	return parseIP(ipstr, true)
}

func ParseIPv6(ipstr string) (net.IP, error) {
	return parseIP(ipstr, false)
}

func parseIP(ipstr string, isv4 bool) (net.IP, error) {
	if ip := net.ParseIP(ipstr); ip == nil || (ip.To4() != nil) != isv4 {
		return nil, fmt.Errorf("invalid ip %s", ipstr)
	} else {
		return ip, nil
	}
}

func IPv4FromUint32(i uint32) net.IP {
	ip := make([]byte, 4)
	binary.BigEndian.PutUint32(ip, i)
	return net.IP(ip)
}

func IPv4ToUint32(ip net.IP) uint32 {
	if ip == nil {
		return 0
	} else {
		return binary.BigEndian.Uint32(ip.To4())
	}
}

func IPv4StrToUint32(ipstr string) (uint32, error) {
	if ipv4, err := parseIP(ipstr, true); err != nil {
		return 0, err
	} else {
		return IPv4ToUint32(ipv4), nil
	}
}

func IPv6FromBigInt(bigint *big.Int) net.IP {
	var ip net.IP
	if bigint != nil {
		ip = net.IP(bigint.Bytes())
	}

	return ip
}

func IPv6ToBigInt(ip net.IP) *big.Int {
	ipv6Int := big.NewInt(0)
	ipv6Int.SetBytes(ip.To16())
	return ipv6Int
}

func IPv6StrToBigInt(ipstr string) (*big.Int, error) {
	if ipv6, err := parseIP(ipstr, false); err != nil {
		return nil, err
	} else {
		return IPv6ToBigInt(ipv6), nil
	}
}
