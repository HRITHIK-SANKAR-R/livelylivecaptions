package network

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const webSocketURL = "ws://localhost:8080/ws" // Placeholder URL

// ManageWebSocket handles the WebSocket connection and data transfer.
func ManageWebSocket(outgoingChan <-chan []byte, incomingChan chan<- string, quitChan <-chan struct{}) {
	interrupt := make(chan struct{})

	u, err := url.Parse(webSocketURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	fmt.Printf("Connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})

	// Receiver goroutine
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			incomingChan <- string(message)
		}
	}()

	// Sender goroutine
	go func() {
		for {
			select {
			case <-done:
				return
			case data := <-outgoingChan:
				err := conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					fmt.Println("Write error:", err)
					return
				}
			}
		}
	}()

	for {
		select {
		case <-quitChan:
			fmt.Println("Closing WebSocket connection.")
			// Cleanly close the connection by sending a close message and then
			// waiting for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("Write close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		case <-interrupt:
			return
		}
	}
}
