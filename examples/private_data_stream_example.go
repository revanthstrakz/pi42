package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/revanthstrakz/pi42"
	"github.com/zishang520/engine.io-client-go/transports"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/socket.io-client-go/socket"
)

func main() {
	PrivateDataStreamExample()
}

func PrivateDataStreamExample() {
	// Load environment variables from .env file
	err := godotenv.Load(filepath.Join(".", ".env"))
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("PI42_API_KEY")
	apiSecret := os.Getenv("PI42_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("API key and secret must be provided. Set PI42_API_KEY and PI42_API_SECRET environment variables.")
	}

	// Create API client
	client := pi42.NewClient(apiKey, apiSecret)

	fmt.Println("=== Authenticated WebSocket Stream Example ===")

	// Step 1: Create a listen key
	fmt.Println("Creating listen key...")
	listenKeyResponse, err := client.UserData.CreateListenKey()
	if err != nil {
		log.Fatalf("Error creating listen key: %v", err)
	}

	listenKey, ok := listenKeyResponse["listenKey"]
	if !ok {
		log.Fatalf("Listen key not found in response: %v", listenKeyResponse)
	}

	fmt.Printf("Listen key obtained: %s\n", listenKey)

	// Create a context with cancellation for clean shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to keep the listen key alive
	keepAliveDone := make(chan struct{})
	go keepAliveListenKey(client, keepAliveDone)

	// Setup WebSocket connection options
	serverUrl := fmt.Sprintf("https://fawss-uds.pi42.com/auth-stream/%s", listenKey)
	fmt.Printf("Connecting to authenticated WebSocket at: %s\n", serverUrl)

	opts := socket.DefaultOptions()
	opts.SetPath("/")
	opts.SetTransports(types.NewSet(
		transports.WebSocket,
		transports.Polling,
	))

	// Create socket manager and socket
	manager := socket.NewManager(serverUrl, opts)
	io := manager.Socket("/", nil)

	// Setup event handlers
	setupConnectionHandlers(io)
	setupUserDataHandlers(io)

	// Run until we get a termination signal
	<-stopChan

	fmt.Println("\nShutting down...")

	// Close the WebSocket connection gracefully
	if io.Connected() {
		io.Disconnect()
	}

	// Stop the keep-alive routine
	close(keepAliveDone)

	// Delete the listen key
	fmt.Println("Deleting listen key...")
	_, err = client.UserData.DeleteListenKey()
	if err != nil {
		log.Printf("Error deleting listen key: %v", err)
	} else {
		fmt.Println("Listen key deleted successfully")
	}

	fmt.Println("=== Authenticated WebSocket Stream Example Completed ===")
}

// keepAliveListenKey periodically updates the listen key to keep it active
func keepAliveListenKey(client *pi42.Client, done <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Minute) // Update every 10 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Updating listen key...")
			result, err := client.UserData.UpdateListenKey()
			if err != nil {
				log.Printf("Error updating listen key: %v", err)
			} else {
				fmt.Printf("Listen key updated: %s\n", result)
			}
		case <-done:
			return
		}
	}
}

// setupConnectionHandlers sets up handlers for connection events
func setupConnectionHandlers(io *socket.Socket) {
	// Connection established
	io.On("connect", func(args ...any) {
		utils.Log().Info("Connected to authenticated WebSocket stream")
	})

	// Connection error
	io.On("connect_error", func(args ...any) {
		utils.Log().Warning("Connection error: %v", args)
	})

	// Disconnection
	io.On("disconnect", func(args ...any) {
		utils.Log().Warning("Disconnected from WebSocket server: %v", args)
	})

	// Handle ping/pong events (for debugging)
	io.On("ping", func(args ...any) {
		utils.Log().Debug("Ping received")
	})

	io.On("pong", func(args ...any) {
		utils.Log().Debug("Pong received")
	})
}

// setupUserDataHandlers sets up handlers for user data events
func setupUserDataHandlers(io *socket.Socket) {
	// New position created
	io.On("newPosition", func(args ...any) {
		fmt.Println("\n📊 New position created:")
		printEventData(args...)
	})

	// Order filled completely
	io.On("orderFilled", func(args ...any) {
		fmt.Println("\n✅ Order filled completely:")
		printEventData(args...)
	})

	// Order partially filled
	io.On("orderPartiallyFilled", func(args ...any) {
		fmt.Println("\n🔄 Order partially filled:")
		printEventData(args...)
	})

	// Order cancelled
	io.On("orderCancelled", func(args ...any) {
		fmt.Println("\n❌ Order cancelled:")
		printEventData(args...)
	})

	// Order failed
	io.On("orderFailed", func(args ...any) {
		fmt.Println("\n❗ Order failed:")
		printEventData(args...)
	})

	// New order (stop order executed)
	io.On("newOrder", func(args ...any) {
		fmt.Println("\n📝 New order (stop order executed):")
		printEventData(args...)
	})

	// Update order (stop limit order executed)
	io.On("updateOrder", func(args ...any) {
		fmt.Println("\n📋 Order updated (stop limit executed):")
		printEventData(args...)
	})

	// Position updated
	io.On("updatePosition", func(args ...any) {
		fmt.Println("\n🔄 Position updated:")
		printEventData(args...)
	})

	// Position closed
	io.On("closePosition", func(args ...any) {
		fmt.Println("\n🚫 Position closed:")
		printEventData(args...)
	})

	// Balance update
	io.On("balanceUpdate", func(args ...any) {
		fmt.Println("\n💰 Balance updated:")
		printEventData(args...)
	})

	// New trade
	io.On("newTrade", func(args ...any) {
		fmt.Println("\n💱 New trade:")
		printEventData(args...)
	})

	// Session expired
	io.On("sessionExpired", func(args ...any) {
		fmt.Println("\n⏰ Session expired:")
		printEventData(args...)
	})
}

// printEventData formats and prints event data for better readability
func printEventData(args ...any) {
	if len(args) == 0 {
		fmt.Println("  [No data]")
		return
	}

	for i, arg := range args {
		switch data := arg.(type) {
		case map[string]interface{}:
			fmt.Printf("  Data %d:\n", i+1)
			for k, v := range data {
				fmt.Printf("    %s: %v\n", k, v)
			}
		default:
			fmt.Printf("  Data %d: %v\n", i+1, arg)
		}
	}
}
