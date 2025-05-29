package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/revanthstrakz/pi42"
)

func LimitCheck() {
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

	// Create a client instance
	client := pi42.NewClient(apiKey, apiSecret)

	// Run the order book and limit order example
	orderBookLimitExample(client)
}

func orderBookLimitExample(client *pi42.Client) {
	// Symbol we want to trade
	symbol := "BTCINR"
	fmt.Printf("\n=== Order Book & Limit Order Example for %s ===\n", symbol)

	// 1. Get the order book depth to analyze current market conditions
	fmt.Println("\n1. Fetching order book depth...")
	depth, err := client.Market.GetDepth(symbol)
	if err != nil {
		log.Fatalf("Error getting order book depth: %v", err)
	}

	// Parse and display order book data
	displayOrderBook(depth)

	// 2. Determine best bid/ask prices and the farthest bid price
	bestBid, bestAsk := getBestPrices(depth)
	farthestBid := getFarthestBidPrice(depth)

	if bestBid == 0 || bestAsk == 0 || farthestBid == 0 {
		log.Fatalf("Could not determine price levels from order book")
	}

	fmt.Printf("Best bid: %.2f\n", bestBid)
	fmt.Printf("Best ask: %.2f\n", bestAsk)
	fmt.Printf("Farthest bid: %.2f\n", farthestBid)

	// 3. Place a limit order at the farthest bid price from the order book
	//    This ensures our order is far from the current market price and unlikely to execute
	limitPrice := farthestBid
	fmt.Printf("\n2. Placing limit buy order at %.2f (farthest bid from order book)...\n", limitPrice)

	// Get contract info for precision
	_, exists := client.ExchangeInfo[symbol]
	if !exists {
		log.Fatalf("Contract information for %s not found", symbol)
	}

	// Use the bullet function which handles precision automatically
	limitOrder, err := client.Order.PlaceOrder(pi42.PlaceOrderParams{
		Symbol:     symbol,
		Side:       pi42.OrderSideBuy,
		Type:       pi42.OrderTypeLimit,
		Price:      limitPrice,
		Quantity:   0.001, // Example quantity, adjust as needed
		PositionID: "limit-order-example",
		ReduceOnly: false, // Not a reduce-only order
		Leverage:   1,     // Leverage, if applicable
		StopPrice:  0,     // No stop price for limit orders
	})

	if err != nil {
		log.Fatalf("Error placing limit order: %v", err)
	}

	fmt.Printf("Limit order placed successfully!\n")
	fmt.Printf("  Order ID: %v\n", limitOrder.ID)
	fmt.Printf("  Client Order ID: %s\n", limitOrder.ClientOrderID)
	fmt.Printf("  Symbol: %s\n", limitOrder.Symbol)
	fmt.Printf("  Price: %.2f\n", limitOrder.Price)
	fmt.Printf("  Quantity: %.8f\n", limitOrder.OrderAmount)

	clientOrderID := limitOrder.ClientOrderID

	// 4. Check open orders to confirm our order is there
	fmt.Println("\n3. Checking open orders...")
	time.Sleep(2 * time.Second) // Wait a moment for the order to be registered
	checkOpenOrders(client, symbol, clientOrderID)

	// 5. Cancel the limit order
	fmt.Printf("\n4. Cancelling order with client order ID: %s...\n", clientOrderID)
	cancelResult, err := client.Order.DeleteOrder(clientOrderID)
	if err != nil {
		log.Fatalf("Error cancelling order: %v", err)
	}

	fmt.Printf("Order cancelled successfully: %v\n", cancelResult)

	// 6. Check open orders again to confirm cancellation
	fmt.Println("\n5. Checking open orders after cancellation...")
	time.Sleep(2 * time.Second) // Wait a moment for cancellation to process
	checkOpenOrders(client, symbol, clientOrderID)

	fmt.Println("\n=== Order Book & Limit Order Example Completed ===")
}

// displayOrderBook parses and displays order book data
func displayOrderBook(depth *pi42.DepthResponse) {
	fmt.Printf("\nOrder Book Data for %s:\n", depth.Data.Symbol)
	fmt.Printf("Event Time: %d\n", depth.Data.EventTime)

	// Display asks (sell orders)
	asks := depth.Data.Asks
	fmt.Println("\nTop 5 asks (sell orders):")
	maxAsks := 5
	if len(asks) < maxAsks {
		maxAsks = len(asks)
	}

	for i := 0; i < maxAsks; i++ {
		askData := asks[i]
		price, _ := strconv.ParseFloat(askData[0], 64)
		quantity, _ := strconv.ParseFloat(askData[1], 64)
		fmt.Printf("  Price: %.2f, Quantity: %.8f\n", price, quantity)
	}

	// Display bids (buy orders)
	bids := depth.Data.Bids
	fmt.Println("\nTop 5 bids (buy orders):")
	maxBids := 5
	if len(bids) < maxBids {
		maxBids = len(bids)
	}

	for i := 0; i < maxBids; i++ {
		bidData := bids[i]
		price, _ := strconv.ParseFloat(bidData[0], 64)
		quantity, _ := strconv.ParseFloat(bidData[1], 64)
		fmt.Printf("  Price: %.2f, Quantity: %.8f\n", price, quantity)
	}
}

// getBestPrices extracts the best bid and ask prices from the order book
func getBestPrices(depth *pi42.DepthResponse) (float64, float64) {
	var bestBid, bestAsk float64

	// Get best ask (lowest sell price)
	if len(depth.Data.Asks) > 0 {
		bestAsk, _ = strconv.ParseFloat(depth.Data.Asks[0][0], 64)
	}

	// Get best bid (highest buy price)
	if len(depth.Data.Bids) > 0 {
		bestBid, _ = strconv.ParseFloat(depth.Data.Bids[0][0], 64)
	}

	return bestBid, bestAsk
}

// getFarthestBidPrice finds the farthest (lowest) bid price in the order book
func getFarthestBidPrice(depth *pi42.DepthResponse) float64 {
	bids := depth.Data.Bids
	if len(bids) == 0 {
		return 0 // No bids available
	}

	// Initialize with the first bid
	farthestPrice, _ := strconv.ParseFloat(bids[0][0], 64)

	// Loop through all bids to find the farthest (lowest) price
	for _, bidData := range bids {
		price, _ := strconv.ParseFloat(bidData[0], 64)
		if price < farthestPrice {
			farthestPrice = price
		}
	}

	// Apply a small additional buffer (0.5% lower) to ensure our order is even farther
	// This helps in case the order book changes slightly between fetching and placing the order
	farthestPrice = farthestPrice * 0.995

	return farthestPrice
}

// checkOpenOrders retrieves and displays open orders
func checkOpenOrders(client *pi42.Client, symbol string, targetOrderID string) {
	openOrders, err := client.Order.GetOpenOrders(pi42.OrderQueryParams{
		Symbol: symbol,
	})

	if err != nil {
		log.Printf("Error getting open orders: %v\n", err)
		return
	}

	fmt.Printf("Found %d open orders for %s\n", len(openOrders), symbol)

	foundTargetOrder := false
	for i, order := range openOrders {
		clientOrderID := order.ClientOrderID

		fmt.Printf("%d. Client Order ID: %s\n", i+1, clientOrderID)
		fmt.Printf("   Symbol: %v\n", order.Symbol)
		fmt.Printf("   Type: %v\n", order.Type)
		fmt.Printf("   Side: %v\n", order.Side)
		fmt.Printf("   Price: %v\n", order.Price)
		fmt.Printf("   Quantity: %v\n", order.OrderAmount)
		fmt.Printf("   Status: %v\n", order.Status)
		fmt.Println()

		if clientOrderID == targetOrderID {
			foundTargetOrder = true
		}
	}

	if targetOrderID != "" {
		if foundTargetOrder {
			fmt.Printf("✓ Target order %s was found in open orders\n", targetOrderID)
		} else {
			fmt.Printf("✗ Target order %s was NOT found in open orders\n", targetOrderID)
		}
	}
}
