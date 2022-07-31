package client

type Configuration struct {
	RemoteAddr string `json:"server"`
	RemotePort int    `json:"port"`
	Secure     bool   `json:"secure"`
}

type Client struct {
	ctrl *Controller
}

func NewClient(config Configuration) (*Client, error) {
	ctrl, err := NewController(config.RemoteAddr, config.RemotePort, config.Secure)
	if err != nil {
		return nil, err
	}

	return &Client{
		ctrl: ctrl,
	}, nil
}

func (c *Client) Release() {
	c.ctrl.Release()
}

func (c *Client) JoinRoom(roomId string) {
	c.ctrl.io.Emit("user:join", roomId)
}
