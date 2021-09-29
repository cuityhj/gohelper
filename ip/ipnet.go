package ip

import (
	"fmt"
	"net"
)

func ParseCIDRv4(cidr string) (*net.IPNet, error) {
	return parseCIDR(cidr, true)
}

func ParseCIDRv6(cidr string) (*net.IPNet, error) {
	return parseCIDR(cidr, false)
}

func parseCIDR(cidr string, isv4 bool) (*net.IPNet, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	} else if (ip.To4() != nil) != isv4 {
		return nil, fmt.Errorf("ipnet %s is invalid ip version", cidr)
	} else {
		return ipnet, nil
	}
}
