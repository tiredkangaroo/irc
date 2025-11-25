package main

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

// Relay relays data between a WebSocket connection and the TCP IRC server.
// It acts as a bridge, forwarding messages in both directions.
// NOTE: this can be exploited to create an open proxy, so this should run
// on the user's machine.
// another NOTE: this function does not close ws conn.
func relay(wc *websocket.Conn, hostport string) error {
	conn, err := net.Dial("tcp", hostport)
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}
	defer conn.Close()

	errc := make(chan error, 2)

	go func() {
		for {
			_, p, err := wc.ReadMessage()
			if err != nil {
				errc <- fmt.Errorf("websocket read error: %v", err)
				break
			}
			_, err = conn.Write(p)
			if err != nil {
				errc <- fmt.Errorf("write to irc error: %v", err)
				break
			}
		}
	}()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				errc <- fmt.Errorf("read from irc error: %v", err)
				break
			}
			err = wc.WriteMessage(websocket.TextMessage, buf[:n])
			if err != nil {
				errc <- fmt.Errorf("websocket write error: %v", err)
				break
			}
		}
	}()

	// wait for an error from either goroutine and return it
	return <-errc
}
