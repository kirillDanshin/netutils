package netutils

import (
	"fmt"
	"log"
	"net"
)

func scan(iface *net.Interface) error {
	var (
		addr  *net.IPNet
		addrs []net.Addr
		err   error
	)

	if addrs, err = iface.Addrs(); err != nil {
		return err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				addr = &net.IPNet{
					IP:   ip4,
					Mask: ipnet.Mask[len(ipnet.Mask)-4:],
				}
				break
			}
		}
	}

	if addr == nil {
		return fmt.Errorf("there's no IP network found")
	}

	if addr.IP[0] == 127 {
		return fmt.Errorf("skipping localhost")
	}

	if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return fmt.Errorf("mask means network is too large")
	}

	log.Printf(
		"Using network range %+v for interface %v",
		addr,
		iface.Name,
	)

	return nil
}
