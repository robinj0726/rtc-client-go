package client

import (
	"crypto/tls"
	"encoding/json"
	"log"

	"github.com/pion/webrtc/v3"
	socketio "github.com/robinj730/rtc-client-go/gosocketio"
	"github.com/robinj730/rtc-client-go/gosocketio/transport"
)

type Controller struct {
	io *socketio.Client
}

func onConnectionHandler(c *socketio.Channel)    { log.Printf("Connected %s\n", c.Id()) }
func onDisconnectionHandler(c *socketio.Channel) { log.Printf("Disconnected %s\n", c.Id()) }

func onJoinedRoom(c *socketio.Channel, data interface{}) {
	j, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	type Payload struct {
		Id         string      `json:"id"`
		IceServers interface{} `json:"iceServers"`
	}

	var obj Payload
	if err := json.Unmarshal(j, &obj); err != nil {
		log.Fatal(err)
	}
	if c.Id() == obj.Id {
		// It's me
		c.Emit("user:rtc:ready", c.Id())
	}
}

func startRtcConnection(c *socketio.Channel, otherSocketId, stunUrl string) {
	log.Printf("ask startRtcConnection from %s\n", otherSocketId)
	webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					stunUrl,
				},
			},
		},
	})
}

func NewController(remoteAddr string, remotePort int, secure bool) (*Controller, error) {
	tr := transport.DefaultWebsocketTransport()
	if secure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	io, err := socketio.Dial(
		socketio.AddrWebsocket(remoteAddr, remotePort, secure),
		tr,
	)
	if err != nil {
		return nil, err
	}

	if err := io.On(socketio.OnConnection, onConnectionHandler); err != nil {
		return nil, err
	}
	if err := io.On(socketio.OnDisconnection, onDisconnectionHandler); err != nil {
		return nil, err
	}

	if err := io.On("user:joined", onJoinedRoom); err != nil {
		return nil, err
	}

	// if err := client.On("user:rtc:start", startRtcConnection); err != nil {
	// 	return nil, err
	// }

	return &Controller{
		io,
	}, nil
}

func (ctrl *Controller) Release() {
	ctrl.io.Close()
}
