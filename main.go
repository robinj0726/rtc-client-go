package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"

	socketio "github.com/mtfelian/golang-socketio"
	"github.com/mtfelian/golang-socketio/transport"
)

func onConnectionHandler(c *socketio.Channel)    { log.Printf("Connected %s\n", c.Id()) }
func onDisconnectionHandler(c *socketio.Channel) { log.Printf("Disconnected %s\n", c.Id()) }

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	tr := transport.DefaultWebsocketTransport()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client, err := socketio.Dial(
		socketio.AddrWebsocket("localhost", 4000, true),
		tr,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.On(socketio.OnConnection, onConnectionHandler); err != nil {
		log.Fatal(err)
	}
	if err := client.On(socketio.OnDisconnection, onDisconnectionHandler); err != nil {
		log.Fatal(err)
	}

	<-interrupt

	client.Close()
}
