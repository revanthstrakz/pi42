package main

import (
	"fmt"
	"time"

	"github.com/revanthstrakz/pi42"

	"github.com/zishang520/engine.io/v2/utils"
)

func DataStreamExample() {
	// Create a new WebSocket client
	client := pi42.NewSocketClient()

	// Define handlers for different data streams
	myHandler := func(data ...any) {
		fmt.Println(".")
	}
	// Initially add some streams
	client.AddStream("btcinr@depth_0.1", "depthUpdate")
	client.AddStream("btcinr@markPrice", "markPriceUpdate")

	// Get channel for the depth event
	depthChannel, exists := client.GetEventChannel("depthUpdate")
	if exists {
		go func() {
			for event := range depthChannel {
				myHandler(event.Data...)
			}
		}()
	}
	// Get channel for the allContractDetails event
	allContractDetailsChannel, exists := client.GetEventChannel("allContractDetails")
	if exists {
		go func() {
			for event := range allContractDetailsChannel {
				myHandler(event.Data...)
			}
		}()
	}

	// Start listening on a specific event channel
	markPriceChannel, exists := client.GetEventChannel("markPriceUpdate")
	if exists {
		go func() {
			for event := range markPriceChannel {
				myHandler(event.Data...)
			}
		}()
	}

	// Start the client in a separate goroutine so we can continue
	// adding and removing streams
	go func() {
		client.Init()
	}()

	// Wait a moment for connection to establish
	time.Sleep(2 * time.Second)

	// Dynamically add a new stream after the client has started
	utils.Log().Info("Adding kline stream...")
	client.AddStream("btcinr@kline_1m", "kline")

	// Get channel for the kline event
	klineChannel, exists := client.GetEventChannel("kline")
	if exists {
		go func() {
			for event := range klineChannel {
				fmt.Printf("Kline event received via channel - Topic: %s, Data: %v\n",
					event.Topic, event.Data)
			}
		}()
	}

	// Add another topic for the same event type (markPriceUpdate)
	utils.Log().Info("Adding another topic for markPriceUpdate event...")
	client.AddStream("ethinr@markPrice", "markPriceUpdate")

	// Now the markPriceUpdate channel will receive events for both topics

	// Wait for some data to arrive
	time.Sleep(5 * time.Second)

	// Example of removing a stream but keeping the event channel active
	utils.Log().Info("Removing btcinr@markPrice but keeping the markPriceUpdate channel...")
	client.RemoveStream("btcinr@markPrice")

	// The markPriceUpdate channel is still active for ethinr@markPrice

	// Wait for another period observing data
	time.Sleep(10 * time.Second)
	select {}
}
