package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"time"
)

type Conn struct{
	conn *websocket.Conn
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{
		conn: conn,
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		log.Printf("Read conn err %v", err)
		return 0, err
	}

	x := copy(b, data)

	return x, err
}

func (c *Conn) Write(b []byte) (n int, err error) {
	err = c.conn.WriteMessage(websocket.BinaryMessage, b)

	return len(b), err
}

func (c *Conn) Close() error {
	log.Printf("WS Conn Close")
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
