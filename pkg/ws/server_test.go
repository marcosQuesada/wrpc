package ws

import (
	log "github.com/sirupsen/logrus"
	"net"
	"testing"
)

func TestConnForwarder(t *testing.T) {

	l := &fakeListener{}
	srv := NewServer(l)
	cin, cout := net.Pipe()
	bin, bout := net.Pipe()
	defer func(){
		cin.Close()
	}()
	go srv.forward(cin, bin)

	payload := []byte("Heeello")
	n, err := cout.Write(payload)
	if err != nil {
		t.Fatalf("Unexpected error writing connection, error %v", err)
	}

	expected := 7
	if n != expected {
		t.Errorf("unexpected write size, expected %d got %d", expected, n)
	}

	rcv := make([]byte, 512)
	n, err = bout.Read(rcv)
	if err != nil {
		t.Fatalf("Unexpected error reading connection, error %v", err)
	}

	if n != expected {
		t.Errorf("unexpected read size, expected %d got %d", expected, n)
	}

	if string(rcv[:n]) != string(payload) {
		log.Errorf("Payloads do not match, expected %s got %s", payload, rcv)
	}
}

type fakeListener struct{
	conn net.Conn
}

func (f *fakeListener) Connect() (net.Conn, error) {
	return f.conn, nil
}
