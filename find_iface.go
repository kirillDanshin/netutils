package netutils

import (
	"log"
	"net"
	"sync"

	"github.com/kirillDanshin/dlog"
)

// FindIfaceWithAddr finds an interface
// with a given address. returns empty string and error
// if not found or error happened.
func FindIfaceWithAddr(addr string) (string, error) {
	ifaceName := ""
	check := func(iface *net.Interface) {
		addrs, err := iface.Addrs()
		if err != nil {
			dlog.F("err: %s", err)
			return
		}
		iface.Addrs()
		for _, ifaceAddr := range addrs {
			dlog.F("ifaceAddr=[%#+v] addr=[%#+v]", ifaceAddr.String(), addr)
			_, cidrnet, err := net.ParseCIDR(ifaceAddr.String())
			if err != nil {
				dlog.F("Error: %s", err)
				return
			}
			myaddr := net.ParseIP(addr)
			dlog.F("contains=[%#+v]", cidrnet.Contains(myaddr))
			if cidrnet.Contains(myaddr) {
				ifaceName = iface.Name
			}
		}
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("err: %s", err)
		return "", err
	}

	err = IfacesWalkSync(ifaces, check)

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
