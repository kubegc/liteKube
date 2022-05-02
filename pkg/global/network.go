package global

import (
	"fmt"
	"math/big"
	"net"
	"strings"

	"k8s.io/klog/v2"
)

var LocalhostIP = net.ParseIP("127.0.0.1")
var LocalHostDNSName = "localhost"
var LocalIPs = QueryIps()
var testUrl = "www.baidu.com:80"
var _, NetworkControllerCIDR, _ = net.ParseCIDR("10.1.1.0/24")

const ReservedNodeToken = "reserverd"

func QueryIps() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		klog.Errorf("fail to get local ips")
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

func RemoveRepeatIps(ips []net.IP) []net.IP {
	new := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		if ip == nil {
			continue
		}

		if !inIps(ip, new) {
			new = append(new, ip)
		}
	}

	return new
}

func inIps(ip net.IP, ips []net.IP) bool {
	for _, m := range ips {
		if net.IP.Equal(ip, m) {
			return true
		}
	}

	return false
}

func GetDefaultServiceIp(clusterIpRange *net.IPNet) net.IP {
	if ip, err := GetIndexedIP(clusterIpRange, 1); err != nil {
		return nil
	} else {
		return ip
	}
}

func GetDefaultClusterDNSIP(clusterIpRange *net.IPNet) net.IP {
	if ip, err := GetIndexedIP(clusterIpRange, 2); err != nil {
		return nil
	} else {
		return ip
	}
}

func GetIndexedIP(subnet *net.IPNet, index int) (net.IP, error) {
	ip := addIPOffset(bigForIP(subnet.IP), index)
	if !subnet.Contains(ip) {
		return nil, fmt.Errorf("can't generate IP with index %d from subnet. subnet too small. subnet: %q", index, subnet)
	}
	return ip, nil
}

func addIPOffset(base *big.Int, offset int) net.IP {
	return net.IP(big.NewInt(0).Add(base, big.NewInt(int64(offset))).Bytes())
}

func bigForIP(ip net.IP) *big.Int {
	b := ip.To4()
	if b == nil {
		b = ip.To16()
	}
	return big.NewInt(0).SetBytes(b)
}

func ExternIp() string {
	conn, err := net.Dial("udp", testUrl)
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	if ip := net.ParseIP(strings.Split(conn.LocalAddr().String(), ":")[0]); ip == nil {
		return "127.0.0.1"
	} else {
		return ip.String()
	}

}
