package network

import (
	"fmt"
	"net"
)

type NetworkManager struct{}

func NewNetworkManager() *NetworkManager {
	return &NetworkManager{}
}

func (n *NetworkManager) GetIPv6Address(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", fmt.Errorf("error getting interface: %v", err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("error getting addresses: %v", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() == nil && !ipnet.IP.IsLinkLocalUnicast() {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no IPv6 address found for interface %s", interfaceName)
}
