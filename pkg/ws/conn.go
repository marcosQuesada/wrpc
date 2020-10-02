package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

// Conn decorates websocket.Conn as net.Conn
type Conn struct {
	conn *websocket.Conn
}

func newConn(conn *websocket.Conn) *Conn {
	return &Conn{
		conn: conn,
	}
}

// Read adapter method
func (c *Conn) Read(b []byte) (n int, err error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		log.Errorf("Read conn err %v", err)
		return 0, err
	}

	x := copy(b, data)

	return x, err
}

// Write adapter method
func (c *Conn) Write(b []byte) (n int, err error) {
	err = c.conn.WriteMessage(websocket.BinaryMessage, b)

	return len(b), err
}

// Close adapter method
func (c *Conn) Close() error {
	log.Printf("WS Conn Close")
	return c.conn.Close()
}

// LocalAddr adapter method
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr adapter method
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline adapter method, not implemented
func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadlinea adapter method
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadLine adapter method
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
