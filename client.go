package gosocketio

import (
	"github.com/fanatic/golang-socketio/transport"
	"strconv"
)

const (
	webSocketProtocol = "ws://"
	webSocketSecureProtocol = "wss://"
	socketioUrl       = "/socket.io/?EIO=3&transport=websocket"
)

/**
Socket.io client representation
*/
type Client struct {
	methods
	Channel
}

/**
Get ws/wss url by host and port
 */
func GetUrl(host string, port int, secure bool) string {
	var prefix string
	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}
	return prefix + host + ":" + strconv.Itoa(port) + socketioUrl
}

/**
connect to host and initialise socket.io protocol

The correct ws protocol url example:
ws://myserver.com/socket.io/?EIO=3&transport=websocket

You can use GetUrlByHost for generating correct url
*/
func Dial(url string, tr transport.Transport) (*Client, error) {
	c := New()

	if err := c.Dial(url, tr); err != nil {
		return nil, err
	}

	return c, nil
}

func New() *Client {
	c := &Client{}
	c.initChannel()
	c.initMethods()
	return c
}

func (c *Client) Dial(url string, tr transport.Transport) error {
	var err error
	c.conn, err = tr.Connect(url)
	if err != nil {
		return err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)

	return nil
}

/**
Close client connection
*/
func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}
