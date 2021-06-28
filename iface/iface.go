package iface

import (
	"fmt"
	"net"
	"strings"
)

func CheckIfaceNameValid(ifname string) error {
	if ifname == "" {
		return nil
	}

	_, err := net.InterfaceByName(ifname)
	return err
}

func GenIfaceAndIp(ifaceAndIp string, isv4 bool) (*net.Interface, net.IP, error) {
	ifnameAndIp := strings.SplitN(ifaceAndIp, "/", 2)
	if len(ifnameAndIp) != 2 {
		return nil, nil, fmt.Errorf("invlaid interface %s", ifaceAndIp)
	}

	iface, err := net.InterfaceByName(ifnameAndIp[0])
	if err != nil {
		return nil, nil, fmt.Errorf("get interface info with ifname %s failed: %s", ifnameAndIp[0], err.Error())
	}

	ip := net.ParseIP(ifnameAndIp[1])
	if ip == nil || (ip.To4() != nil) != isv4 {
		return nil, nil, fmt.Errorf("invalid ip %s", ifnameAndIp[1])
	}

	return iface, ip, checkIpBelongsToIface(iface, ip)
}

func checkIpBelongsToIface(iface *net.Interface, ip net.IP) error {
	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok == false {
			continue
		}

		if ipnet.IP.Equal(ip) {
			return nil
		}
	}

	return fmt.Errorf("network interface %s no addr %s", iface.Name, ip.String())
}

func GetGlobalUnicastIPs(isv4 bool) ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if isv4 == false && (iface.Flags&net.FlagMulticast) != net.FlagMulticast {
				continue
			}

			ipnet, ok := addr.(*net.IPNet)
			if ok == false {
				continue
			}

			if ip := ipnet.IP; (ip.To4() != nil) == isv4 && ip.IsGlobalUnicast() {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips, nil
}
