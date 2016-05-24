package whisper

import (
	"fmt"
	"net"
)

// ExternalIP looks up an the first available external IP address.
func ExternalIP() (string, *Error) {

	// Get addresses for the interface
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", &Error{
			text: fmt.Sprintf("Could not get interface addresses: %s", err.Error()),
			code: 96,
		}
	}

	// Go through each address to find a an IPv4
	for _, addr := range addrs {

		var ip net.IP

		switch val := addr.(type) {
		case *net.IPNet:
			ip = val.IP
		case *net.IPAddr:
			ip = val.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue // ignore loopback and nil addresses
		}

		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}

		return ip.String(), nil
	}

	return "", &Error{
		text: "Are you connected to the network?!",
		code: 95,
	}

}
