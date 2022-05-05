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
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/Litekube/network-controller/certs"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/contant"
	"github.com/gorilla/websocket"
	"github.com/songgao/water"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"time"

	"fmt"
)

type Client struct {
	stopCh chan struct{}
	// config
	cfg config.ClientConfig
	// interface
	iface *water.Interface
	// ip addr
	ip      net.IP
	toIface chan []byte
	ws      *websocket.Conn
	data    chan *Data
	state   int
	// store route des ip
	routes          []string
	ClientTLSConfig config.TLSConfig
}

var net_gateway, net_nic string

func NewClient(cfg config.ClientConfig) *Client {

	if cfg.MTU != 0 {
		MTU = cfg.MTU
	}

	client := &Client{
		cfg:     cfg,
		iface:   nil,
		ip:      nil,
		toIface: make(chan []byte, 100),
		ws:      nil,
		data:    make(chan *Data, 100),
		state:   0,
		routes:  make([]string, 0, 1024),
		ClientTLSConfig: config.TLSConfig{
			CAFile:         cfg.CAFile,
			CAKeyFile:      filepath.Join(cfg.NetworkCertDir, contant.CAKeyFile),
			ClientCertFile: cfg.ClientCertFile,
			ClientKeyFile:  cfg.ClientKeyFile,
		},
	}
	return client
}

func (client *Client) Run() error {

	go client.cleanUp()

	iface, err := newTun("")
	if err != nil {
		return err
	}
	client.iface = iface

	net_gateway, net_nic, err = GetNetGateway()
	logger.Debugf("Net gateway:%+v, nic:%+v", net_gateway, net_nic)
	if err != nil {
		logger.Error(err)
		return err
	}

	srvDest := client.cfg.ServerAddr + "/32"
	addRoute(srvDest, net_gateway, net_nic)
	client.routes = append(client.routes, srvDest)

	// build ws connect to network server
	srvAdr := fmt.Sprintf("%s:%d", client.cfg.ServerAddr, client.cfg.Port)
	u := url.URL{Scheme: "wss", Host: srvAdr, Path: "/ws"}
	header := http.Header{}
	header.Set(contant.NodeTokenKey, client.cfg.Token)
	logger.Debugf("Connecting to %+v", u.String())

	// continue to try to connect every 2s until success
	// fix here, conenct immediately，then 0.5s
	// ticker := time.NewTicker(3 * time.Second)
	var connection *websocket.Conn
	logger.Infof("client try to connect %+v", u.String())

	// load ca & client key/cert
	pool, err := certs.LoadCertPool(client.ClientTLSConfig.CAFile)
	if err != nil {
		logger.Error(err)
		return err
	}
	tlsCert, err := tls.LoadX509KeyPair(client.ClientTLSConfig.ClientCertFile, client.ClientTLSConfig.ClientKeyFile)
	if err != nil {
		logger.Error(err)
		return err
	}
	for ok := true; ok; ok = (connection == nil) {
		// support tls
		dialer := websocket.Dialer{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{tlsCert},
			},
		}
		connection, _, err = dialer.Dial(u.String(), header)
		//connection, _, err = websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			logger.Infof("Dial: %+v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	client.ws = connection
	defer connection.Close()

	// init state
	client.state = contant.STATE_INIT

	client.ws.SetReadLimit(maxMessageSize)
	client.ws.SetReadDeadline(time.Now().Add(pongWait))
	// client send ping, receive pong
	// SetPongHandler sets the handler for pong messages received from the peer.
	client.ws.SetPongHandler(func(string) error {
		client.ws.SetReadDeadline(time.Now().Add(pongWait))
		logger.Debug("Pong received")
		return nil
	})

	go client.writePump()

	// Initialize connection with master
	client.data <- &Data{
		ConnectionState: contant.STATE_CONNECT,
	}

	// read
	for {
		messageType, r, err := connection.ReadMessage()
		if err != nil {
			logger.Error(err)
			delRoute("0.0.0.0/1")
			delRoute("128.0.0.0/1")
			for _, dest := range client.routes {
				delRoute(dest)
			}
			break
		} else {
			if messageType == websocket.TextMessage {
				client.dispatcher(r)
			}
		}
	}
	return errors.New("Not expected to exit")
}

func (client *Client) dispatcher(p []byte) {
	logger.Debugf("Dispatcher client state: %+v", client.state)
	switch client.state {
	case contant.STATE_INIT:
		logger.Debug("STATE_INIT")
		var message Data
		if err := json.Unmarshal(p, &message); err != nil {
			client.ws.Close()
			close(client.data)
			logger.Panic(err)
		}
		if message.ConnectionState == contant.STATE_CONNECT {

			ipStr := string(message.Payload)
			ip, subnet, _ := net.ParseCIDR(ipStr)
			setTunIP(client.iface, ip, subnet)
			if client.cfg.RedirectGateway {
				err := redirectGateway(client.iface.Name(), tun_peer.String())
				if err != nil {
					logger.Errorf("Redirect gateway error: %+v", err.Error())
				}
			}

			client.state = contant.STATE_CONNECTED
			client.handleInterface()
		}
	case contant.STATE_CONNECTED:
		// write data to local interface channel
		client.toIface <- p
	}
}

func (client *Client) handleInterface() {
	// network packet to interface
	go func() {
		for {
			hp := <-client.toIface
			_, err := client.iface.Write(hp)
			if err != nil {
				logger.Errorf("handleInterface write iface err:%+v", err)
				return
			}
			logger.Debug("Write to interface")
		}
	}()

	// interface to network packet
	go func() {
		packet := make([]byte, contant.IFACE_BUFSIZE)
		for {
			plen, err := client.iface.Read(packet)
			if err != nil {
				logger.Errorf("handleInterface read iface err: %+v", err)
				break
			}
			client.data <- &Data{
				ConnectionState: contant.STATE_CONNECTED,
				Payload:         packet[:plen],
			}
		}
	}()
}

func (client *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.ws.Close()
	}()

	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				client.write(websocket.CloseMessage, &Data{})
				return
			}
			if err := client.write(websocket.TextMessage, message); err != nil {
				logger.Errorf("client.write err: %+v", err)
			}
		case <-ticker.C:
			// heartbeat 30s
			if err := client.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				logger.Error("Send ping error", err)
			}
		}
	}
}

func (client *Client) write(mt int, message *Data) error {

	if message.ConnectionState == contant.STATE_CONNECTED {
		return client.ws.WriteMessage(mt, message.Payload)
	} else {
		s, err := json.Marshal(message)
		if err != nil {
			logger.Panic(err)
		}
		return client.ws.WriteMessage(mt, s)
	}

}

// client exit gracefully
func (client *Client) cleanUp() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logger.Info("Cleaning Up")
	// redirectGateway = true
	delRoute("0.0.0.0/1")
	delRoute("128.0.0.0/1")
	for _, dest := range client.routes {
		delRoute(dest)
	}

	os.Exit(0)
}
