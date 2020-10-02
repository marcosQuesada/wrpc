package bufconn

import (
	"fmt"
	"net"
	"sync"
)

// Listener implements a net.Listener that creates local, buffered net.Conns
// via its Accept and Connect method.
type Listener struct {
	mu   sync.Mutex
	ch   chan net.Conn
	done chan struct{}
}

var errClosed = fmt.Errorf("closed")

// Listen returns a Listener that can only be contacted by its own Dialers and
// creates buffered connections between the two.
func Listen() *Listener {
	return &Listener{ch: make(chan net.Conn), done: make(chan struct{})}
}

// Accept blocks until Dial is called, then returns a net.Conn for the server
// half of the connection.
func (l *Listener) Accept() (net.Conn, error) {
	select {
	case <-l.done:
		return nil, errClosed
	case c := <-l.ch:
		return c, nil
	}
}

// Close stops the listener.
func (l *Listener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	select {
	case <-l.done:
		// Already closed.
		break
	default:
		close(l.done)
	}
	return nil
}

// Addr reports the address of the listener.
func (l *Listener) Addr() net.Addr { return addr{} }

// Connect attach one pipe side to listener server
func (l *Listener) Connect() (net.Conn, error) {
	inBound, conn := net.Pipe()
	select {
	case <-l.done:
		return inBound, errClosed
	case l.ch <- conn:
		return inBound, nil
	}
}

type addr struct{}

func (addr) Network() string { return "bufconn" }
func (addr) String() string  { return "bufconn" }
