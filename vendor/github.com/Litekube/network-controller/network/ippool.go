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
	"errors"
	"fmt"
	"github.com/Litekube/network-controller/sqlite"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
)

/*
assign unique ip for client
*/

type NetworkIpPool struct {
	subnet *net.IPNet
	pool   [127]int32 // map cache
}

var poolFull = errors.New("IP Pool Full")

// bind existed ip or get an empty ip
func (p *NetworkIpPool) next(bindIp string) (*net.IPNet, error) {

	// assign ip+mask
	ipnet := &net.IPNet{
		make([]byte, 4),
		make([]byte, 4),
	}
	copy([]byte(ipnet.IP), []byte(p.subnet.IP))
	copy([]byte(ipnet.Mask), []byte(p.subnet.Mask))

	if len(bindIp) != 0 {
		ip3, err := strconv.Atoi(strings.Split(bindIp, ".")[3])
		if atomic.LoadInt32(&p.pool[ip3]) == 0 {
			// unnormal
			return nil, errors.New("internal error: conflict between cache and db")
		}
		if err != nil {
			return nil, errors.New(fmt.Sprint("bind existed ip err: %+v", err.Error()))
		}
		ipnet.IP[3] = byte(ip3)
		return ipnet, nil
	}

	found := false
	var i int
	// server take x.1 & x.2, begin from 3
	for i = 3; i < 255; i += 2 {
		// CAS sync
		if atomic.CompareAndSwapInt32(&p.pool[i], 0, 1) {
			found = true
			break
		}
	}

	// find db with LRU
	if !found {
		nm := sqlite.NetworkMgr{}
		item, _ := nm.QueryLogestIdle()
		if item != nil {
			// no need to set cache=0, just delete old item from sqlite
			found = true
			i, _ = strconv.Atoi(strings.Split(bindIp, ".")[3])
			res, err := nm.DeleteById(item.Id)
			if !res || err != nil {
				logger.Errorf("fail to delete idle item: %+v", err)
				return nil, errors.New(fmt.Sprint("fail to delete idle item: %+v", err.Error()))
			}
		}
	}

	if !found {
		return nil, poolFull
	}

	ipnet.IP[3] = byte(i) // found=true
	return ipnet, nil
}

// release ip
func (p *NetworkIpPool) release(ip net.IP) {
	defer func() {
		// recover only work in defer part
		// if normal, return nil
		// if panic, return panic err and recover normal,continue to execute
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	logger.Infof("releasing ip: %+v", ip)
	i := ip[3]
	p.pool[i] = 0
}

func (p *NetworkIpPool) releaseByTag(tag int) {
	defer func() {
		// recover only work in defer part
		// if normal, return nil
		// if panic, return panic err and recover normal,continue to execute
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	//logger.Infof("releasing ip: %+v", tag)
	p.pool[tag] = 0
}
