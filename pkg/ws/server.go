package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
)

var defaultBufSize = 32 * 1024

var upgrader = websocket.Upgrader{} // use default options

type Listener interface {
	Connect() (net.Conn, error)
}

type server struct {
	listener Listener
}

func NewServer(l Listener) *server {
	return &server{
		listener: l,
	}
}

func (f *server) Handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}

	var conn net.Conn = newConn(c)
	defer conn.Close()

	inBound, err := f.listener.Connect()
	if err != nil {
		log.Print("dial error:", err)
		return
	}

	go f.readPump(inBound, conn)
	f.writePump(inBound, conn)
}

func (f *server) readPump(inBound, conn net.Conn) {
	for {
		data := make([]byte, defaultBufSize)
		n, err := conn.Read(data)
		if err != nil {
			log.Println("Error ReadMessage:", err)
			_ = inBound.Close()
			break
		}

		_, err = inBound.Write(data[:n])
		if err != nil {
			log.Println("inbound write error:", err)
			break
		}
	}
}

func (f *server) writePump(inBound, conn net.Conn) {
	for {
		rsp := make([]byte, defaultBufSize)
		n, err := inBound.Read(rsp)
		if err != nil {
			log.Println("readAll:", err)
			break
		}

		_, err = conn.Write(rsp[:n])
		if err != nil {
			log.Println("piped write:", err)
			break
		}
	}
}
