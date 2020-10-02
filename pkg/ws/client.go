package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/url"
)

// NewClient builds net.Conn from websocket connection
func NewClient(u url.URL) (net.Conn, error) {
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("Error dialing, %v \n", err)
		return nil, err
	}

	return newConn(c), nil
}
