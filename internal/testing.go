package internal

import (
	"net"
	"testing"
)

// TestConn is a helper function for returning a client and server
// net.Conn connected to each other.
// NOTE: taken from https://github.com/hashicorp/go-plugin/blob/cfdf485783602a2ca85502dedebf441be7bcbc8d/testing.go#L36
func TestConn(t testing.TB) (net.Conn, net.Conn) {
	// Listen to any local port. This listener will be closed
	// after a single connection is established.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Start a goroutine to accept our client connection
	var serverConn net.Conn
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		defer l.Close()
		var err error
		serverConn, err = l.Accept()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	// Connect to the server
	clientConn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Wait for the server side to acknowledge it has connected
	<-doneCh

	return clientConn, serverConn
}
