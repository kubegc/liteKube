package global

import (
	"fmt"
	"net"
)

var LocalhostIP = net.ParseIP("127.0.0.1")
var LocalHostDNSName = "localhost"
var LocalIPs = QueryIps()

const ReservedNodeToken = "reserverd"

func QueryIps() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return []net.IP{LocalhostIP}
	}

	ips := make([]net.IP, 0, 2)
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}

	return append(ips, LocalhostIP)
}
