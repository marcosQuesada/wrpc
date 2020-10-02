package bufconn

import (
	"errors"
	"testing"
	"time"
)

func TestAcceptorBlocksUntilConnectionIsPassed(t *testing.T) {
	l := NewListener()

	done := make(chan struct{})
	go func() {
		_, err := l.Accept()
		if err != nil && !errors.Is(errClosed, err) {
			t.Fatalf("unexpeccted error accepting, error %v", err)
		}

		close(done)
	}()

	_, err := l.Connect()
	if err != nil {
		t.Fatalf("unexpeccted error on connect, error %v", err)
	}
	timeout := time.NewTimer(time.Second)
	select{
	case <- timeout.C:
		t.Fatal("Unexpected timeout error")
	case <- done:
	}
}

func TestAcceptorGetsUnlockedOnCancel(t *testing.T) {
	l := NewListener()

	done := make(chan struct{})
	go func() {
		_, err := l.Accept()
		if err != nil && !errors.Is(errClosed, err) {
			t.Fatalf("unexpeccted error accepting, error %v", err)
		}

		close(done)
	}()

	_ = l.Close()
	timeout := time.NewTimer(time.Second)
	select{
	case <- timeout.C:
		t.Fatal("Unexpected timeout error")
	case <- done:
	}
}