package pi42

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zishang520/engine.io-client-go/transports"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/socket.io-client-go/socket"
)

// EventData represents data received from a WebSocket event
type EventData struct {
	// The event name (like depthUpdate, markPriceUpdate)
	Event types.EventName
	// The specific topic for this event (like btcinr@depth_0.1)
	Topic string
	// The data received from the WebSocket
	Data []any
}

// SocketClient is a client for WebSocket connections
type SocketClient struct {
	// Socket client instance
	io *socket.Socket
	// Manager instance for handling connections
	manager *socket.Manager
	// List of events to subscribe to
	events []types.EventName
	// List of topics to subscribe to
	topics []string
	// Channels for events, mapped by event name
	eventChannels map[types.EventName]chan EventData
	// Mutex for thread-safe access to channels
	channelMutex sync.RWMutex
}

// NewSocketClient creates a new WebSocket client
func NewSocketClient() *SocketClient {
	ec := make(map[types.EventName]chan EventData)
	for _, event := range []types.EventName{
		"depthUpdate",
		"markPriceUpdate",
		"kline",
		"aggTrade",
		"24hrTicker",
		"marketInfo",
		"tickerArr",
		"markPriceArr",
		"allContractDetails",
	} {
		ec[event] = make(chan EventData) // Buffered channel for each event
	}
	return &SocketClient{
		events: []types.EventName{
			"depthUpdate",
			"markPriceUpdate",
			"kline",
			"aggTrade",
			"24hrTicker",
			"marketInfo",
			"tickerArr",
			"markPriceArr",
			"allContractDetails",
		},
		topics:        []string{},
		eventChannels: ec,
	}
}

// AddStream adds a new topic and corresponding event handler
func (sc *SocketClient) AddStream(topic string, event types.EventName) {
	// Check if topic already exists
	for _, t := range sc.topics {
		if t == topic {
			return // Topic already exists
		}
	}

	sc.topics = append(sc.topics, topic)

	// If already connected, subscribe to the new topic immediately
	if sc.io != nil && sc.io.Connected() {
		sc.io.Emit("subscribe", map[string][]string{
			"params": {topic},
		})
	}
}

// RemoveStream removes a specific topic from the subscription list
func (sc *SocketClient) RemoveStream(topic string) {
	// Find and remove the topic from the list
	for i, t := range sc.topics {
		if t == topic {
			// Remove the topic using slice manipulation
			sc.topics = append(sc.topics[:i], sc.topics[i+1:]...)

			// If already connected, unsubscribe from the topic immediately
			if sc.io != nil && sc.io.Connected() {
				sc.io.Emit("unsubscribe", map[string][]string{
					"params": {topic},
				})
				utils.Log().Info("Unsubscribed from topic: %s", topic)
			}
			return
		}
	}

	utils.Log().Warning("Topic not found for removal: %s", topic)
}

// GetEventChannel returns a channel for a specific event
func (sc *SocketClient) GetEventChannel(event types.EventName) (chan EventData, bool) {
	sc.channelMutex.RLock()
	defer sc.channelMutex.RUnlock()

	ch, exists := sc.eventChannels[event]
	return ch, exists
}

func (sc *SocketClient) Init() {
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	opts := socket.DefaultOptions()
	opts.SetTransports(types.NewSet(transports.Polling, transports.WebSocket))

	// Updated server URL
	manager := socket.NewManager("https://fawss.pi42.com/", opts)
	sc.manager = manager

	// Listening to manager events
	sc.manager.On("error", func(errs ...any) {
		utils.Log().Warning("Manager Error: %v", errs)
	})

	sc.manager.On("ping", func(...any) {
		utils.Log().Warning("Manager Ping")
	})

	sc.manager.On("reconnect", func(...any) {
		utils.Log().Warning("Manager Reconnected")
	})

	sc.manager.On("reconnect_attempt", func(...any) {
		utils.Log().Warning("Manager Reconnect Attempt")
	})

	sc.manager.On("reconnect_error", func(errs ...any) {
		utils.Log().Warning("Manager Reconnect Error: %v", errs)
	})

	sc.manager.On("reconnect_failed", func(errs ...any) {
		utils.Log().Warning("Manager Reconnect Failed: %v", errs)
	})

	// Using default namespace
	io := sc.manager.Socket("/", opts)
	sc.io = io

	// Print detailed socket information for debugging
	utils.Log().Info("Socket object initialized: %v", io)
	utils.Log().Info("Socket ID: %v", io.Id())
	utils.Log().Info("Socket connected: %v", io.Connected())

	sc.io.On("connect", func(args ...any) {
		utils.Log().Info("Connected to WebSocket server, ID: %v", io.Id())
		utils.Log().Info("Connection state: %v", io.Connected())

		// Subscribe to topics after connection is established
		subscribeToTopics(sc)
	})

	sc.io.On("connect_error", func(args ...any) {
		utils.Log().Warning("Connection error: %v", args)

		// Attempt to reconnect after error
		if !io.Connected() {
			utils.Log().Info("Attempting to reconnect...")
			io.Connect()
		}
	})

	sc.io.On("disconnect", func(args ...any) {
		utils.Log().Warning("Disconnected from WebSocket server: %+v", args)
	})

	// Wait for termination signal
	<-sigChan
	utils.Log().Info("Shutting down...")

	// Clean disconnect
	if sc.io.Connected() {
		sc.io.Disconnect()
	}
}

// Helper function to subscribe to configured topics
func subscribeToTopics(sc *SocketClient) {
	if len(sc.topics) == 0 {
		utils.Log().Info("No topics to subscribe to")
		return
	}

	utils.Log().Info("Subscribing to topics: %v", sc.topics)

	// Subscribe to each topic by emitting the subscribe event
	sc.io.Emit("subscribe", map[string][]string{
		"params": sc.topics,
	})

	// Add an acknowledgment callback for the subscription
	sc.io.EmitWithAck("subscribe", func(ack ...any) {
		utils.Log().Info("Subscription acknowledgment: %v", ack)
	}, map[string][]string{
		"params": sc.topics,
	})

	// Setup event handlers with debug output
	setupEventHandlers(sc)
}

// Function to set up all event handlers
func setupEventHandlers(sc *SocketClient) {
	// Setup a single handler for each event type
	for _, event := range sc.events {
		// Create a handler that can determine which topic triggered the event
		eventHandler := createChannelEventHandler(sc, event)
		setupEventHandler(sc.io, event, eventHandler)
	}
}

func createChannelEventHandler(sc *SocketClient, event types.EventName) func(...any) {
	eventchannel, exists := sc.GetEventChannel(event)
	if exists {
		return func(data ...any) {
			select {
			case eventchannel <- EventData{
				Event: event,
				Data:  data,
			}:
				// Message sent successfully
			default:
				// Channel buffer is full, log a warning
				utils.Log().Warning("Channel buffer full for event %s; dropping message", event)
			}
		}
	}
	return func(data ...any) {
		utils.Log().Warning("Event channel not found for event: %s", event)
	}
}

func setupEventHandler(io *socket.Socket, event types.EventName, function func(...any)) {
	io.On(event, function)
}
