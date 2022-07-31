package client

import (
	"crypto/tls"
	"fmt"
	"log"

	socketio "github.com/robinj730/rtc-client-go/gosocketio"
	"github.com/robinj730/rtc-client-go/gosocketio/transport"
)

type Connection struct {
	conn *socketio.Client
}

func onConnectionHandler(c *socketio.Channel)    { log.Printf("Connected %s\n", c.Id()) }
func onDisconnectionHandler(c *socketio.Channel) { log.Printf("Disconnected %s\n", c.Id()) }

func NewClient(remoteAddr string, remotePort int, secure bool) (*Connection, error) {
	tr := transport.DefaultWebsocketTransport()
	if secure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client, err := socketio.Dial(
		socketio.AddrWebsocket(remoteAddr, remotePort, secure),
		tr,
	)
	if err != nil {
		return nil, err
	}

	if err := client.On(socketio.OnConnection, onConnectionHandler); err != nil {
		return nil, err
	}
	if err := client.On(socketio.OnDisconnection, onDisconnectionHandler); err != nil {
		return nil, err
	}

	return &Connection{
		conn: client,
	}, nil
}

func (c *Connection) Release() {
	c.conn.Close()
}

func (c *Connection) JoinRoom(roomId string) {
	c.conn.On("user:joined", func(c *socketio.Channel, data string) {
		fmt.Printf("#%v\n", data)
	})
	c.conn.Emit("user:join", roomId)
}
