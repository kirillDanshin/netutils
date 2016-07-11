package netutils

import (
	"log"
	"net"
	"sync"

	"github.com/kirillDanshin/dlog"
	"github.com/kirillDanshin/myutils"
)

// FindIfaceWithAddr finds an interface
// with a given address. returns empty string and error
// if not found or error happened.
func FindIfaceWithAddr(addr string, withCaller ...bool) (string, error) {
	var wCaller bool
	if len(withCaller) > 0 {
		wCaller = withCaller[0]
	}
	ifaceName := ""
	check := func(iface *net.Interface) {
		addrs, err := iface.Addrs()
		if err != nil {
			if wCaller {
				dlog.F("%s: [\n\terr: %s\n]", myutils.Slice(dlog.GetCaller(2))[0].(*dlog.Caller).String(), err)
			} else {
				dlog.F("err: %s", err)
			}
			return
		}
		iface.Addrs()
		for _, ifaceAddr := range addrs {
			if wCaller {
				dlog.F(
					"%s: [\n\tifaceAddr=[%#+v]\n\taddr=[%#+v]\n]",
					myutils.Slice(dlog.GetCaller(2))[0].(*dlog.Caller).String(),
					ifaceAddr.String(),
					addr,
				)
			} else {
				dlog.F("ifaceAddr=[%#+v] addr=[%#+v]", ifaceAddr.String(), addr)
			}
			_, cidrnet, err := net.ParseCIDR(ifaceAddr.String())
			if err != nil {
				if wCaller {
					dlog.F("%s: [\n\tError: %s\n]", myutils.Slice(dlog.GetCaller(2))[0].(*dlog.Caller).String(), err)
				} else {
					dlog.F("Error: %s", err)
				}
				return
			}
			myaddr := net.ParseIP(addr)
			if wCaller {
				dlog.F("%s: [\n\tcontains=[%#+v]\n]",
					myutils.Slice(dlog.GetCaller(2))[0].(*dlog.Caller).String(),
					cidrnet.Contains(myaddr),
				)
			} else {
				dlog.F("contains=[%#+v]", cidrnet.Contains(myaddr))
			}
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
