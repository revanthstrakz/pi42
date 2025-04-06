package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	origin := "http://localhost/"
	serverURL := "https://fawss.pi42.com/socket.io/"

	ws, err := websocket.Dial(serverURL, "", origin)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer ws.Close()

	fmt.Println("Connected to WebSocket server")

	namespaceReady := make(chan bool, 1)
	stop := make(chan struct{})

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ws.Read(buf)
			if err != nil {
				log.Println("Read error:", err)
				// close(stop)
				// return
			}

			// The Socket.IO protocol may send multiple messages at once
			// e.g., "2" or "42[...]2" in a single read
			messages := strings.Split(string(buf[:n]), "\x1e") // "\x1e" is the message separator used in newer Socket.IO (EIO4+)
			for _, msg := range messages {
				msg = strings.TrimSpace(msg)
				if msg == "" {
					continue
				}
				fmt.Println("Received:", msg)

				switch {
				case msg == "2":
					// Server is pinging us
					ws.Write([]byte("3"))
					fmt.Println("Sent pong")
				case msg == "3":
					fmt.Println("Received pong")
				case msg == "40":
					fmt.Println("Namespace ready")
					select {
					case namespaceReady <- true:
					default:
					}
				case strings.HasPrefix(msg, "42"):
					fmt.Println("Event:", msg)
				case strings.HasPrefix(msg, "0"):
					fmt.Println("Handshake complete")
					ws.Write([]byte("40")) // trigger default namespace
				default:
					fmt.Println("Unhandled:", msg)
				}
			}
		}
	}()

	// Subscribe only once when namespace is ready
	go func() {
		select {
		case <-namespaceReady:
			time.Sleep(500 * time.Millisecond)
			subscribe := `42["subscribe",{"params":["btcinr@kline _ 1m"]}]`
			_, err := ws.Write([]byte(subscribe))
			if err != nil {
				log.Println("Subscribe error:", err)
				return
			}
			fmt.Println("Sent subscribe message")
		case <-time.After(10 * time.Second):
			log.Println("Timed out waiting for namespace readiness")
		}
	}()

	// Start ping loop only after namespace is ready
	go func() {
		select {
		case <-namespaceReady:
			fmt.Println("Starting ping loop...")
		case <-time.After(10 * time.Second):
			log.Println("Timeout waiting for ping start")
			return
		}

		ticker := time.NewTicker(170 * time.Second) // server expects every 180s
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := ws.Write([]byte("2"))
				if err != nil {
					log.Println("Ping failed:", err)
					return
				}
				fmt.Println("Sent ping")
			case <-stop:
				return
			}
		}
	}()

	// Block forever (until EOF)
	select {
	case <-stop:
		log.Println("Connection closed, exiting.")
	}
}
