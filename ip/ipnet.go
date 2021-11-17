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
	} else if ipnet.IP.Equal(ip) == false {
		return nil, fmt.Errorf("ipnet %s no match %s", cidr, ipnet.String())
	} else if (ip.To4() != nil) != isv4 {
		return nil, fmt.Errorf("ipnet %s no match ip version", cidr)
	} else {
		return ipnet, nil
	}
}
