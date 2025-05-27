package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/revanthstrakz/pi42"
	"github.com/zishang520/engine.io/v2/utils"
)

func wsex() {
	// Setup signal handling for graceful shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// Create a new WebSocket client
	client := pi42.NewSocketClient()

	// Subscribe to specific topics before starting the client
	fmt.Println("Subscribing to initial streams...")
	client.AddStream("btcinr@depth_0.1", "depthUpdate")
	client.AddStream("btcinr@markPrice", "markPriceUpdate")
	client.AddStream("btcinr@kline_1m", "kline")

	// Setup handlers for different event types
	setupEventHandlers(client)

	// Start the client in a separate goroutine
	go func() {
		fmt.Println("Starting WebSocket client...")
		client.Init()
	}()

	// Wait a moment for connection to establish
	time.Sleep(2 * time.Second)

	// Add more streams after client has started
	fmt.Println("Adding additional streams...")
	client.AddStream("ethinr@markPrice", "markPriceUpdate")
	client.AddStream("ethinr@kline_1m", "kline")

	// Dynamic management example: adding and removing streams with delay
	go manageStreams(client)

	// Wait for termination signal
	<-signalChannel
	fmt.Println("Shutting down...")
}

// setupEventHandlers configures handlers for different event types
func setupEventHandlers(client *pi42.SocketClient) {
	// Handle depth updates
	depthChannel, exists := client.GetEventChannel("depthUpdate")
	if exists {
		go func() {
			for event := range depthChannel {
				// Print only first update from each batch to avoid console spam
				if len(event.Data) > 0 {
					fmt.Printf("Depth update received from %s\n", event.Topic)
				}
			}
		}()
	}

	// Handle mark price updates
	markPriceChannel, exists := client.GetEventChannel("markPriceUpdate")
	if exists {
		go func() {
			for event := range markPriceChannel {
				fmt.Printf("Mark price update from %s: %v\n", event.Topic, event.Data)
			}
		}()
	}

	// Handle kline/candlestick updates
	klineChannel, exists := client.GetEventChannel("kline")
	if exists {
		go func() {
			for event := range klineChannel {
				fmt.Printf("Kline from %s: Interval=%v\n", event.Topic, event.Data)
			}
		}()
	}

	// Handle ticker updates
	tickerChannel, exists := client.GetEventChannel("24hrTicker")
	if exists {
		go func() {
			for event := range tickerChannel {
				fmt.Printf("Ticker update from %s\n", event.Topic)
			}
		}()
	}
}

// manageStreams demonstrates how to add and remove streams dynamically
func manageStreams(client *pi42.SocketClient) {
	// Wait 10 seconds before adding more streams
	time.Sleep(10 * time.Second)

	utils.Log().Info("Adding ticker stream for BTCINR...")
	client.AddStream("btcinr@ticker", "24hrTicker")

	// Wait 5 seconds before removing a stream
	time.Sleep(5 * time.Second)

	utils.Log().Info("Removing btcinr@markPrice stream...")
	client.RemoveStream("btcinr@markPrice")

	// Note that the markPriceUpdate channel remains active and will continue
	// to receive updates from other markPrice streams (like ethinr@markPrice)

	// Wait 5 seconds before adding the stream back
	time.Sleep(5 * time.Second)

	utils.Log().Info("Adding back btcinr@markPrice stream...")
	client.AddStream("btcinr@markPrice", "markPriceUpdate")

	utils.Log().Info("Stream management completed")
}
