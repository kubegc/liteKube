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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/sqlite"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait) / 2
	maxMessageSize = 1024 * 1024 //1MB
)

type connection struct {
	id        int
	ws        *websocket.Conn
	server    *NetworkServer
	data      chan *Data
	state     int // STATE_INIT / STATE_CONNECTED
	ipAddress *net.IPNet
	token     string
	bindIp    string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

var maxId int = 0

func NewConnection(ws *websocket.Conn, server *NetworkServer, token string) (*connection, error) {
	if ws == nil {
		panic("ws cannot be nil")
	}
	if server == nil {
		panic("server cannot be nil")
	}

	nm := sqlite.NetworkMgr{}
	item, _ := nm.QueryByToken(token)
	bindIp := ""
	if item == nil {
		return nil, errors.New(fmt.Sprintf("invalid token %+v", token))
	}
	//item != nil
	if item.State != contant.STATE_IDLE {
		return nil, errors.New(fmt.Sprintf("token %+v already connected", token))
	}
	if len(item.BindIp) != 0 {
		bindIp = item.BindIp
	}
	// auto inc
	maxId++
	data := make(chan *Data)
	c := &connection{maxId, ws, server, data, contant.STATE_INIT, nil, token, bindIp}

	// fix server gen token, no need insert now
	if len(bindIp) == 0 {
		//nm.Insert(sqlite.NetworkMgr{
		//	Token:  token,
		//	State:  STATE_INIT,
		//	BindIp: "",
		//})
		nm.UpdateStateByToken(contant.STATE_INIT, token)
	}

	go c.writePump()
	go c.readPump()
	logger.Debug("New connection created")

	return c, nil
}

func (c *connection) readPump() {
	defer func() {
		c.server.unregister <- c
		c.ws.Close()
	}()

	// If a message exceeds the limit, the connection sends a close message to the peer
	c.ws.SetReadLimit(maxMessageSize)

	// heartbeat
	// SetPingHandler sets the handler for ping messages received from the peer, default pong
	// server receive ping, send pong
	c.ws.SetPingHandler(func(string) error {
		logger.Debug("Ping received")
		// WriteControl writes a control message with the given 10s deadline.
		if err := c.ws.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
			logger.Errorf("Send ping error:%+v", err)
		}
		return nil
	})

	// continue to read
	for {
		messageType, r, err := c.ws.ReadMessage()
		if err == io.EOF {
			c.cleanUp()
			break
		} else if err != nil {
			logger.Info(err)
			c.cleanUp()
			break
		} else {
			if messageType == websocket.TextMessage {
				c.dispatcher(r)
			}
		}
	}
}

func (c *connection) writePump() {

	defer func() {
		c.ws.Close()
	}()

	// continue to write
	for {
		if c != nil {
			select {
			case message, ok := <-c.data:
				// Thread can be still active after close connection
				if message != nil {
					logger.Debugf("writePump data len: %+v", len(message.Payload))
					if !ok {
						c.write(websocket.CloseMessage, &Data{})
						return
					}
					if err := c.write(websocket.TextMessage, message); err != nil {
						logger.Errorf("writePump err:%+v", err)
					}
				} else {
					break
				}
			}
		} else {
			break
		}
	}
}

func (c *connection) write(mt int, message *Data) error {

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))

	if message.ConnectionState == contant.STATE_CONNECTED {
		// write payload
		err := c.ws.WriteMessage(mt, message.Payload)
		if err != nil {
			return err
		}
	} else {
		// write payload+connectionState
		s, err := json.Marshal(message)
		if err != nil {
			logger.Panic(err)
			return err
		}

		err = c.ws.WriteMessage(mt, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *connection) dispatcher(p []byte) {
	logger.Debug("Dispatcher connection %+v state: ", c.ipAddress, c.state)
	switch c.state {
	case contant.STATE_INIT:
		logger.Debug("STATE_INIT")
		var message Data
		if err := json.Unmarshal(p, &message); err != nil {
			logger.Panic(err)
		}
		// receive client connect message
		if message.ConnectionState == contant.STATE_CONNECT {
			d := &Data{}
			d.ConnectionState = contant.STATE_CONNECT

			cltIP, err := c.server.ippool.next(c.bindIp)
			if err != nil {
				c.cleanUp()
				logger.Error(err)
			}
			if len(c.bindIp) == 0 {
				nm := sqlite.NetworkMgr{}
				nm.UpdateIpByToken(cltIP.IP.String(), c.token)
			}

			logger.Infof("get next IP from ippool %+v", cltIP)
			d.Payload = []byte(cltIP.String())

			// change connection parameter
			c.ipAddress = cltIP
			c.state = contant.STATE_CONNECTED
			// after connected, register
			c.server.register <- c
			c.data <- d
		}
	case contant.STATE_CONNECTED:
		// if connected, write to channel(tun0)
		logger.Debug("STATE_CONNECTED")
		c.server.toIface <- p
	}
}

func (c *connection) cleanUp() {
	// client close connection
	c.server.unregister <- c
	c.ws.Close()
}
