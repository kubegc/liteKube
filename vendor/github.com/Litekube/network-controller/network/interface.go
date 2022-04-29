/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: wanna <wananzjx@163.com>
 *
 */
package network

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/songgao/water"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
)

var invalidAddr = errors.New("Invalid device ip address")

var tun_peer net.IP

// create tun interface
func newTun(name string) (iface *water.Interface, err error) {

	iface, err = water.New(water.Config{})
	if err != nil {
		return nil, err
	}
	logger.Infof("created interface %v", iface.Name())

	// exec script
	scmd := fmt.Sprintf("ip link set dev %s up mtu %d qlen 100", iface.Name(), MTU)
	cmd := exec.Command("bash", "-c", scmd)
	logger.Infof("exec command: ip %s", scmd)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return iface, nil
}

func setTunIP(iface *water.Interface, ip net.IP, subnet *net.IPNet) (err error) {
	ip = ip.To4()
	logger.Infof("parse network addr ip:%+v,subnet:%+v", ip, subnet)
	// 10.1.1.1 valid & 10.1.1.2 invalid
	if ip[3]%2 == 0 {
		return invalidAddr
	}

	// 4 bytes for ipv4
	peer := net.IP(make([]byte, 4))
	copy([]byte(peer), []byte(ip))
	peer[3]++
	// 10.1.1.1+1
	tun_peer = peer

	// assign ip for tun0
	// ip addr add dev tun0 local 10.1.1.1 peer 10.1.1.2
	scmd := fmt.Sprintf("ip addr add dev %s local %s peer %s", iface.Name(), ip, peer)
	cmd := exec.Command("bash", "-c", scmd)
	logger.Infof("exec command: ip %+v", scmd)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// des 10.1.1.0/24，gateway to 10.1.1.2
	// ip route add 10.1.1.0/24 via 10.1.1.2 dev tun0
	scmd = fmt.Sprintf("ip route add %s via %s dev %s", subnet, peer, iface.Name())
	cmd = exec.Command("bash", "-c", scmd)
	logger.Infof("exec command: ip %+v", scmd)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

// return net gateway (default route) and nic
func GetNetGateway() (gw, dev string, err error) {

	file, err := os.Open("/proc/net/route")
	//file, err := os.Open("./test.txt")
	if err != nil {
		return "", "", err
	}

	defer file.Close()
	rd := bufio.NewReader(file)

	// divide into four parts: 006BA8C0 - 192.168.107.0
	s2byte := func(s string) byte {
		// each part 8bit, ox
		b, _ := strconv.ParseUint(s, 16, 8)
		return byte(b)
	}

	// skip title row
	rd.ReadLine()
	for {
		line, isPrefix, err := rd.ReadLine()

		if err != nil {
			logger.Error(err.Error())
			if err == io.EOF {
				return "", "", errors.New("No default gateway found")
			}
			return "", "", err
		}
		// isPrefix=true indicate Line Too Long
		if isPrefix {
			return "", "", errors.New("Line Too Long!")
		}

		/*
			line example:
			Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRT
			ens160	00000000	026BA8C0	0003	0		0	0		00000000	0	0		0
		*/
		buf := bytes.NewBuffer(line)
		scanner := bufio.NewScanner(buf)
		// split text by space
		scanner.Split(bufio.ScanWords)
		tokens := make([]string, 0, 8)

		for scanner.Scan() {
			tokens = append(tokens, scanner.Text())
		}
		iface := tokens[0]
		dest := tokens[1]
		gw := tokens[2]
		mask := tokens[7]

		// find default interface: dest & mast = 0.0.0.0
		if bytes.Equal([]byte(dest), []byte("00000000")) &&
			bytes.Equal([]byte(mask), []byte("00000000")) {
			// divide into four parts
			a := s2byte(gw[6:8])
			b := s2byte(gw[4:6])
			c := s2byte(gw[2:4])
			d := s2byte(gw[0:2])
			// order is reversed 006BA8C0 - 192.168.107.0
			ip := net.IPv4(a, b, c, d)

			return ip.String(), iface, nil
		}
	}
}

// add route
func addRoute(dest, nextHop, iface string) {

	// ip -4 r a 101.43.253.110/32 via 192.168.107.2 dev ens160
	scmd := fmt.Sprintf("ip -4 r a %s via %s dev %s", dest, nextHop, iface)
	cmd := exec.Command("bash", "-c", scmd)
	logger.Infof("exec command: %+v", scmd)
	err := cmd.Run()

	if err != nil {
		logger.Warning(err.Error())
	}
}

// delete route
func delRoute(dest string) {
	scmd := fmt.Sprintf("ip -4 route del %s", dest)
	cmd := exec.Command("bash", "-c", scmd)
	logger.Infof("exec command: %s", scmd)
	err := cmd.Run()

	if err != nil {
		logger.Warning(err.Error())
	}
}

// redirect default gateway
func redirectGateway(iface, gw string) error {
	subnets := []string{"0.0.0.0/1", "128.0.0.0/1"}
	logger.Info("Redirecting Gateway")
	for _, subnet := range subnets {
		scmd := fmt.Sprintf("ip -4 route add %s via %s dev %s", subnet, gw, iface)
		cmd := exec.Command("bash", "-c", scmd)
		logger.Infof("exec command: ip %s", scmd)
		err := cmd.Run()

		if err != nil {
			return err
		}
	}
	return nil
}
