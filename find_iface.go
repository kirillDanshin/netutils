package netutils

import (
	"net"
	"sync"
)

// FindIfaceWithAddr finds an interface
// with a given address. returns empty string and error
// if not found or error happened.
func FindIfaceWithAddr(addr string) (string, error) {
	ifaceName := ""
	check := func(iface *net.Interface) {
		addrs, err := iface.Addrs()
		if err != nil {
			return
		}
		for _, ifaceAddr := range addrs {
			if ifaceAddr.String() == addr {
				ifaceName = iface.Name
			}
		}
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	if len(ifaces) >= 8 {
		err = IfacesWalk(ifaces, check)
	} else {
		err = IfacesWalkSync(ifaces, check)
	}

	return ifaceName, err
}

// IfacesWalk walks through given net.Interface slice
// in goroutines
func IfacesWalk(ifaces []net.Interface, f func(*net.Interface)) error {
	var wg sync.WaitGroup
	for _, iface := range ifaces {
		wg.Add(1)
		go func(iface net.Interface) {
			defer wg.Done()
			f(&iface)
		}(iface)
	}

	// wait for all interfaces' scan to complete.
	wg.Wait()

	return nil
}

// IfacesWalkSync walks through given net.Interface slice
// in a sync loop
func IfacesWalkSync(ifaces []net.Interface, f func(*net.Interface)) error {
	for _, iface := range ifaces {
		f(&iface)
	}
	return nil
}
