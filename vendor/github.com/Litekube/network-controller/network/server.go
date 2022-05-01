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
	"fmt"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/grpc/grpc_server"
	"github.com/Litekube/network-controller/sqlite"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type NetworkServer struct {
	// config
	cfg config.ServerConfig
	// interface
	iface *water.Interface
	// subnet
	ipnet *net.IPNet
	// IP Pool
	ippool *NetworkIpPool
	// client peers, key is the mac address, value is a HopPeer record
	// Registered clients clientip-connection
	clients map[string]*connection
	// Register requests
	register chan *connection
	// Unregister requests
	unregister chan *connection
	//outData        *Data
	//inData         chan *Data
	toIface          chan []byte
	wg               sync.WaitGroup
	unRegisterCh     chan string
	idleCheckTimer   *time.Ticker
	networkTLSConfig config.TLSConfig
}

var networkServer *NetworkServer

//func GetNetworkServer() *NetworkServer {
//	return networkServer
//}

func NewServer(cfg config.ServerConfig) *NetworkServer {

	if cfg.MTU != 0 {
		MTU = cfg.MTU
	}

	networkServer = &NetworkServer{
		cfg:            cfg,
		iface:          nil,
		ipnet:          nil,
		ippool:         &NetworkIpPool{},
		clients:        make(map[string]*connection),
		register:       make(chan *connection),
		unregister:     make(chan *connection),
		toIface:        make(chan []byte, 100),
		wg:             sync.WaitGroup{},
		unRegisterCh:   nil,
		idleCheckTimer: time.NewTicker(contant.IdleTokenCheckDuration),
		networkTLSConfig: config.TLSConfig{
			CAFile:         cfg.NetworkCAFile,
			CAKeyFile:      cfg.NetworkCAKeyFile,
			ServerCertFile: cfg.NetworkServerCertFile,
			ServerKeyFile:  cfg.NetworkServerKeyFile,
			ClientCertFile: filepath.Join(cfg.NetworkCertDir, contant.ClientCertFile),
			ClientKeyFile:  filepath.Join(cfg.NetworkCertDir, contant.ClientKeyFile),
		},
	}
	return networkServer
}

func (server *NetworkServer) Run() error {

	unRegisterCh := make(chan string, 8)
	networkServer.unRegisterCh = unRegisterCh
	gServer := grpc_server.NewGrpcServer(server.cfg, unRegisterCh)
	go gServer.StartGrpcServerTcp()
	go gServer.StartBootstrapServerTcp()

	//utils.CreateDir(server.cfg.NetworkCertDir)
	//err := certs.CheckNetworkCertConfig(networkServer.networkTLSConfig)
	//if err != nil {
	//	return err
	//}
	// sync cache with db
	networkServer.wg = sync.WaitGroup{}
	networkServer.wg.Add(1)
	go networkServer.initSyncBindIpWithDb()
	go networkServer.handleGrpcUnRegister()

	iface, err := newTun("")
	if err != nil {
		return err
	}
	networkServer.iface = iface

	// networkaddr = 10.1.1.1/24
	ip, subnet, err := net.ParseCIDR(server.cfg.NetworkAddr)
	err = setTunIP(iface, ip, subnet)
	if err != nil {
		return err
	}
	networkServer.ipnet = &net.IPNet{ip, subnet.Mask}
	networkServer.ippool.subnet = subnet

	go networkServer.cleanUp()
	go networkServer.run()

	networkServer.handleInterface()

	// http handle for client to connect
	router := mux.NewRouter()
	router.HandleFunc("/ws", networkServer.serveWs)
	addr := fmt.Sprintf(":%d", networkServer.cfg.Port)

	// wait for cache&db sync
	networkServer.wg.Wait()
	logger.Infof("server ready to ListenAndServe at %+v", addr)
	//err = http.ListenAndServe(addr, router)
	err = http.ListenAndServeTLS(addr, networkServer.networkTLSConfig.ServerCertFile, networkServer.networkTLSConfig.ServerKeyFile, router)
	if err != nil {
		logger.Panicf("ListenAndServe: %+v", err.Error())
	}
	return nil
}

func (server *NetworkServer) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	token := r.Header.Get(contant.NodeTokenKey)
	logger.Infof("reqeust from token: %+v", token)
	// client http to ws
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	// invalid token, close ws conn
	_, err = NewConnection(ws, server, token)
	if err != nil {
		logger.Warning(err)
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
	}
}

func (server *NetworkServer) run() {
	for {
		select {
		case c := <-server.register:
			// add to clients
			logger.Infof("Connection registered: %+v", c.ipAddress.IP.String())
			server.clients[c.ipAddress.IP.String()] = c
			nm := sqlite.NetworkMgr{}
			nm.UpdateStateByToken(contant.STATE_CONNECTED, c.token)
			break

		case c := <-server.unregister:
			// remove from clients
			// close connection data channel
			// release client ip
			clientIP := c.ipAddress.IP.String()
			_, ok := server.clients[clientIP]
			if ok {
				delete(server.clients, clientIP)
				close(c.data)
				if c.ipAddress != nil {
					// unregister for stable ip
					// server.ippool.release(c.ipAddress.IP)
					nm := sqlite.NetworkMgr{}
					nm.UpdateStateByToken(contant.STATE_IDLE, c.token)
				}
				logger.Infof("unregister Connection: %+v, current active clients number: %+v", c.ipAddress.IP, len(server.clients))
			}
			break
		}
	}
}

func (server *NetworkServer) handleInterface() {
	// network packet to interface
	go func() {
		for {
			hp := <-server.toIface
			logger.Debug("Write to interface")
			_, err := server.iface.Write(hp)
			if err != nil {
				logger.Error(err.Error())
				return
			}

		}
	}()

	// interface to network packet
	go func() {
		packet := make([]byte, contant.IFACE_BUFSIZE)
		for {
			plen, err := server.iface.Read(packet)
			if err != nil {
				logger.Error(err)
				break
			}
			header, _ := ipv4.ParseHeader(packet[:plen])
			logger.Debugf("Try sending: %+v", header)
			clientIP := header.Dst.String()
			client, ok := server.clients[clientIP]
			if ok {
				// config file "interconnection=false" not allowed connection between clients
				if !server.cfg.Interconnection {
					if server.isConnectionBetweenClients(header) {
						logger.Infof("Drop connection betwenn %+v and %+v", header.Src, header.Dst)
						continue
					}
				}

				logger.Debugf("Sending to client: %+v", client.ipAddress)
				client.data <- &Data{
					ConnectionState: contant.STATE_CONNECTED,
					Payload:         packet[:plen],
				}

			} else {
				logger.Warningf("Client not found: %+v", clientIP)
			}
		}
	}()
}

func (server *NetworkServer) isConnectionBetweenClients(header *ipv4.Header) bool {

	// srcip!= server ip & desip=one client ip
	if header.Src.String() != header.Dst.String() && header.Src.String() != server.ipnet.IP.String() && server.ippool.subnet.Contains(header.Dst) {
		return true
	}
	return false
}

// server exit gracefully
func (server *NetworkServer) cleanUp() {

	c := make(chan os.Signal, 1)
	// watch ctrl+c or kill pid
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logger.Debug("clean up")

	// update all connected state in sqlite
	nm := sqlite.NetworkMgr{}
	_, err := nm.UpdateAllState()
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("update all connected state")

	// close all client connection
	for key, client := range server.clients {
		client.ws.Close()
		delete(server.clients, key)
	}
	close(server.unregister)

	// code zero indicates success
	os.Exit(0)
}

func (server *NetworkServer) initSyncBindIpWithDb() error {
	defer server.wg.Done()
	nm := sqlite.NetworkMgr{}
	ipList, err := nm.QueryAll()
	if err != nil {
		return err
	}
	logger.Debugf("ipList: %+v", ipList)
	for _, ip := range ipList {
		// register token only, not connect yet
		if len(ip) != 0 {
			tag, _ := strconv.Atoi(strings.Split(ip, ".")[3])
			// no Concurrency
			networkServer.ippool.pool[tag] = 1
		}
	}
	// ignore exsit err, guarantee for reserverd
	nm.Insert(sqlite.NetworkMgr{
		Token:  contant.ReservedToken,
		State:  -1,
		BindIp: "",
	})
	return nil
}

func (server *NetworkServer) handleGrpcUnRegister() error {
	logger.Infof("start handle unregister ip channel")
	for {
		select {
		case ip := <-server.unRegisterCh:
			logger.Infof("receive ip: %+v", ip)
			// close connection
			c, ok := server.clients[ip]
			// may close before unRegister grpc
			if ok {
				delete(server.clients, ip)
				close(c.data)
				c.ws.Close()
			}
			// release ip
			tag, _ := strconv.Atoi(strings.Split(ip, ".")[3])
			server.ippool.releaseByTag(tag)
		case <-server.idleCheckTimer.C:
			tm := sqlite.TokenMgr{}
			tm.DeleteExpireToken()
		}
	}
}
