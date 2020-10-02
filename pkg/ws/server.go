package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io"
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

// NewServer creates a websocket server with piped connections
func NewServer(l Listener) *server {
	return &server{
		listener: l,
	}
}

// Handler handles websocket connection, once established inbound and outbound traffic is forwarded between connections
func (f *server) Handler(w http.ResponseWriter, r *http.Request) {
	defer log.Info("Handler goes down")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade error:", err)
		return
	}

	var conn net.Conn = newConn(c)
	inBound, err := f.listener.Connect()
	if err != nil {
		log.Errorf("dial error:", err)
		return
	}

	go f.forward(conn, inBound)
	f.forward(inBound, conn)
}

func (f *server) forward(inBound, outbound net.Conn) {
	defer log.Info("Forwarder goes down")
	defer outbound.Close()
	for {
		rsp := make([]byte, defaultBufSize)
		n, err := inBound.Read(rsp)
		if err != nil {
			if _, ok := (err).(*websocket.CloseError);!ok && err != io.ErrClosedPipe {
				log.Errorf("Inbound read error: %v", err)
			}
			break
		}

		_, err = outbound.Write(rsp[:n])
		if err != nil {
			log.Errorf("piped write error: %v", err)
			break
		}
	}
}
