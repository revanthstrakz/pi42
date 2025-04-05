package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"path/filepath"
	"pi42"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(filepath.Join("..", ".env"))
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("PI42_API_KEY")
	apiSecret := os.Getenv("PI42_API_SECRET")

	// Create a client instance
	client := pi42.NewClient(apiKey, apiSecret)

	// Run examples
	publicAPIExamples(client)
	authenticatedAPIExamples(client)
	websocketExample(client)
}

func publicAPIExamples(client *pi42.Client) {
	fmt.Println("\n=== Public API Examples ===")

	// Get exchange info
	exchangeInfo, err := client.Exchange.ExchangeInfo("")
	if err != nil {
		fmt.Printf("Error getting exchange info: %v\n", err)
	} else {
		contracts, ok := exchangeInfo["contracts"].([]interface{})
		if ok {
			fmt.Printf("Exchange Info: Found %d contracts\n", len(contracts))
		} else {
			fmt.Println("Exchange Info: Could not parse contracts")
		}
	}

	// Get ticker data for BTC
	ticker, err := client.Market.GetTicker24hr("BTCINR")
	if err != nil {
		fmt.Printf("Error getting ticker: %v\n", err)
		ticker, err = client.Market.GetTicker24hr("BTCINR")
		if err != nil {
			fmt.Printf("Also failed with BTCINR: %v\n", err)
		}
	}

	if err == nil {
		data, ok := ticker["data"].(map[string]interface{})
		if ok {
			fmt.Printf("BTC 24hr Ticker: Last price = %v\n", data["c"])
		} else {
			fmt.Println("BTC 24hr Ticker: Could not parse data")
		}
	}

	// Get klines data
	klines, err := client.Market.GetKlines(pi42.KlinesParams{
		Pair:     "BTCINR",
		Interval: "1h",
		Limit:    5,
	})
	if err != nil {
		fmt.Printf("Error getting klines: %v\n", err)
	} else {
		fmt.Printf("BTCINR Klines: Retrieved %d hourly candles\n", len(klines))
		// Print the first candle to verify data
		if len(klines) > 0 {
			fmt.Printf("First candle: Open=%v, Close=%v\n", klines[0]["open"], klines[0]["close"])
		}
	}
}

func authenticatedAPIExamples(client *pi42.Client) {
	fmt.Println("\n=== Authenticated API Examples ===")

	if client.APIKey == "" || client.APISecret == "" {
		fmt.Println("Skipped - No API Keys")
		return
	}

	// Get wallet details
	wallet, err := client.Wallet.FuturesWalletDetails("INR")
	if err != nil {
		fmt.Printf("Error getting wallet details: %v\n", err)
	} else {
		fmt.Printf("Futures Wallet: Available balance = %v INR\n", wallet["withdrawableBalance"])
	}

	wallet, err = client.Wallet.FundingWalletDetails("INR")
	if err != nil {
		fmt.Printf("Error getting funding wallet details: %v\n", err)
	} else {
		fmt.Printf("Funding Wallet: Available balance = %v INR\n", wallet["withdrawableBalance"])
	}

	// Get open orders
	orders, err := client.Order.GetOpenOrders(pi42.OrderQueryParams{})
	if err != nil {
		fmt.Printf("Error getting open orders: %v\n", err)
	} else {
		fmt.Printf("Open Orders: Found %d open orders\n", len(orders))
	}

	// Get open positions
	positions, err := client.Position.GetPositions("OPEN", pi42.PositionQueryParams{})
	if err != nil {
		fmt.Printf("Error getting positions: %v\n", err)
	} else {
		fmt.Printf("Open Positions: Found %d open positions\n", len(positions))
	}
}

func websocketExample(client *pi42.Client) {
	fmt.Println("\n=== WebSocket Example ===")

	// Track received messages
	receivedMessages := 0

	// Create a function to handle incoming data
	handleTicker := func(data map[string]interface{}) {
		receivedMessages++
		fmt.Printf("Received ticker update for %v: Price = %v\n", data["s"], data["c"])
	}

	// Register the callback
	client.WebSocket.On("24hrTicker", handleTicker)

	// Connect to WebSocket and subscribe to BTCINR ticker
	fmt.Println("Connecting to WebSocket and subscribing to BTCINR ticker...")
	err := client.WebSocket.ConnectPublic([]string{"btcinr@ticker"})
	if err != nil {
		fmt.Printf("WebSocket connection error: %v\n", err)
		return
	}

	// Wait for some data
	fmt.Println("Waiting for data (will exit after 30 seconds or when data is received)...")

	// Wait up to 30 seconds for data
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			if receivedMessages == 0 {
				fmt.Println("No ticker updates received within timeout period")
			} else {
				fmt.Printf("Received %d ticker updates\n", receivedMessages)
			}
			client.WebSocket.Close()
			return
		case <-ticker.C:
			if receivedMessages > 0 {
				fmt.Printf("Received %d ticker updates\n", receivedMessages)
				client.WebSocket.Close()
				return
			}
		}
	}
}
