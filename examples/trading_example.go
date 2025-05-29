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

func tradeex() {
	// Load environment variables from .env file
	err := godotenv.Load(filepath.Join("..", ".env"))
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

	// Check wallet balance
	fmt.Println("\n=== Wallet Information ===")
	checkWalletBalance(client)

	// Get information about trading contracts
	fmt.Println("\n=== Contract Information ===")
	symbol := "BTCINR" // Change to your desired trading pair
	// contractInfo := getContractInfo(client, symbol)

	// Place different types of orders
	fmt.Println("\n=== Order Examples ===")
	placeOrders(client, symbol)

	// Get active positions
	fmt.Println("\n=== Position Information ===")
	activePositions, _ := client.Position.GetPositions(pi42.PositionStatusOpen, pi42.PositionQueryParams{})

	// Find a position for the specified symbol
	var positionID string
	for _, position := range activePositions {
		if position.ContractPair == symbol {
			positionID = position.PositionID
			fmt.Printf("Found position for %s: ID=%s, Size=%f, Entry Price=%f\n",
				symbol, positionID, position.PositionSize, position.EntryPrice)
			break
		}
	}

	// If a position exists, demonstrate position management
	if positionID != "" {
		fmt.Println("\n=== Position Management ===")
		managePosition(client, positionID)
	}

	// Get order history
	fmt.Println("\n=== Order History ===")
	getOrderHistory(client, symbol)

	// Demonstrate order cancellation (Uncomment if you want to test)
	// fmt.Println("\n=== Order Cancellation ===")
	// cancelOrders(client)

	fmt.Println("\n=== Trading Example Completed ===")
}

// checkWalletBalance gets and displays wallet balance information
func checkWalletBalance(client *pi42.Client) {
	// Get futures wallet details
	futuresWallet, err := client.Wallet.FuturesWalletDetails("INR")
	if err != nil {
		log.Printf("Error getting futures wallet details: %v\n", err)
		return
	}

	fmt.Printf("Futures Wallet:\n")
	fmt.Printf("  Available Balance: %v INR\n", futuresWallet.WithdrawableBalance)
	fmt.Printf("  Total Balance: %v INR\n", futuresWallet.WalletBalance)

	// Get funding wallet details
	fundingWallet, err := client.Wallet.FundingWalletDetails("INR")
	if err != nil {
		log.Printf("Error getting funding wallet details: %v\n", err)
		return
	}

	fmt.Printf("Funding Wallet:\n")
	fmt.Printf("  Available Balance: %v INR\n", fundingWallet.WithdrawableBalance)
	fmt.Printf("  Total Balance: %v INR\n", fundingWallet.WalletBalance)
}

// getContractInfo gets and displays information about a specific trading contract
func getContractInfo(client *pi42.Client, symbol string) pi42.ContractInfo {
	// Check if we have the information in the exchange info cache
	contractInfo, exists := client.ExchangeInfo[symbol]
	if !exists {
		log.Printf("Contract information for %s not found in cache\n", symbol)
		return pi42.ContractInfo{}
	}

	fmt.Printf("Contract Information for %s:\n", symbol)
	fmt.Printf("  Base Asset: %s\n", contractInfo.BaseAsset)
	fmt.Printf("  Quote Asset: %s\n", contractInfo.QuoteAsset)
	fmt.Printf("  Price Precision: %d\n", contractInfo.PricePrecision)
	fmt.Printf("  Quantity Precision: %d\n", contractInfo.QuantityPrecision)
	fmt.Printf("  Min Quantity: %f\n", contractInfo.MinQuantity)
	fmt.Printf("  Max Quantity: %f\n", contractInfo.MaxQuantity)
	fmt.Printf("  Market Min Quantity: %f\n", contractInfo.MarketMinQuantity)
	fmt.Printf("  Market Max Quantity: %f\n", contractInfo.MarketMaxQuantity)
	fmt.Printf("  Max Leverage: %f\n", contractInfo.MaxLeverage)
	fmt.Printf("  Supported Order Types: %v\n", contractInfo.OrderTypes)
	fmt.Printf("  Margin Assets: %v\n", contractInfo.MarginAssets)

	// Update leverage preference (optional)
	result, err := client.Exchange.UpdateLeverage(5, symbol)
	if err != nil {
		log.Printf("Error updating leverage: %v\n", err)
	} else {
		fmt.Printf("  Leverage updated successfully: %v\n", result)
	}

	return contractInfo
}

// placeOrders demonstrates how to place different types of orders
func placeOrders(client *pi42.Client, symbol string) {
	// Get current market price
	ticker, err := client.Market.GetTicker24hr(symbol)
	if err != nil {
		log.Printf("Error getting ticker: %v\n", err)
		return
	}

	var currentPrice float64
	if data, ok := ticker["data"].(map[string]interface{}); ok {
		if lastPrice, ok := data["c"].(string); ok {
			// Parse the string to float64
			currentPrice, err = strconv.ParseFloat(lastPrice, 64)
			if err != nil {
				log.Printf("Error parsing price: %v\n", err)
				return
			}
		}
	}

	if currentPrice == 0 {
		log.Println("Could not determine current price")
		return
	}

	fmt.Printf("Current price of %s: %f\n", symbol, currentPrice)

	// Place a small market order (buy)
	fmt.Println("Placing market buy order...")
	marketBuyOrder, err := client.Order.Bullet(pi42.BulletParams{
		Symbol:     symbol,
		Side:       "BUY",
		OrderType:  "MARKET",
		Count:      0.5, // Use a small multiplier for testing
		ReduceOnly: false,
	})

	if err != nil {
		log.Printf("Error placing market buy order: %v\n", err)
	} else {
		fmt.Printf("Market buy order placed successfully!\n")
		fmt.Printf("  Order ID: %v\n", marketBuyOrder.ID)
		fmt.Printf("  Quantity: %f\n", marketBuyOrder.OrderAmount)
		fmt.Printf("  Locked Margin: %f %s\n",
			marketBuyOrder.LockedMargin, marketBuyOrder.MarginAsset)
	}

	// Wait a moment before placing next order
	time.Sleep(1 * time.Second)

	// Place a limit sell order above current price
	limitPrice := currentPrice * 1.02 // 2% above current price
	fmt.Printf("Placing limit sell order at %f...\n", limitPrice)
	limitSellOrder, err := client.Order.Bullet(pi42.BulletParams{
		Symbol:     symbol,
		Side:       "SELL",
		OrderType:  "LIMIT",
		Price:      limitPrice,
		Count:      0.5,  // Use a small multiplier for testing
		ReduceOnly: true, // This will only reduce an existing position
	})

	if err != nil {
		log.Printf("Error placing limit sell order: %v\n", err)
	} else {
		fmt.Printf("Limit sell order placed successfully!\n")
		fmt.Printf("  Order ID: %v\n", limitSellOrder.ID)
		fmt.Printf("  Price: %f\n", limitSellOrder.Price)
		fmt.Printf("  Quantity: %f\n", limitSellOrder.OrderAmount)
	}
}

// managePosition demonstrates how to manage an existing position
func managePosition(client *pi42.Client, positionID string) {
	// Get position details
	positionDetails, err := client.Position.GetPosition(positionID)
	if err != nil {
		log.Printf("Error getting position details: %v\n", err)
		return
	}

	fmt.Printf("Position Details: %v\n", positionDetails)

	// Add margin to the position
	fmt.Println("Adding margin to position...")
	addMarginResult, err := client.Order.AddMargin(positionID, 10) // Add 10 units of margin asset
	if err != nil {
		log.Printf("Error adding margin: %v\n", err)
	} else {
		fmt.Printf("Added margin successfully: %v\n", addMarginResult)
	}

	// Wait a moment
	time.Sleep(1 * time.Second)

	// Reduce margin from the position
	fmt.Println("Reducing margin from position...")
	reduceMarginResult, err := client.Order.ReduceMargin(positionID, 5) // Reduce by 5 units
	if err != nil {
		log.Printf("Error reducing margin: %v\n", err)
	} else {
		fmt.Printf("Reduced margin successfully: %v\n", reduceMarginResult)
	}

	// Place a take profit order for this position
	fmt.Println("Placing take profit order for position...")
	takeProfitOrder, err := client.Order.PlaceOrder(pi42.PlaceOrderParams{
		Symbol:      positionDetails.BaseAsset + positionDetails.QuoteAsset,
		Side:        "SELL", // Assuming it's a long position
		Type:        "LIMIT",
		Quantity:    positionDetails.PositionSize,
		Price:       positionDetails.EntryPrice * 1.05, // 5% profit
		MarginAsset: positionDetails.MarginAsset,
		ReduceOnly:  true,
		PositionID:  positionID,
		PlaceType:   "POSITION", // Specify this is for a specific position
	})

	if err != nil {
		log.Printf("Error placing take profit order: %v\n", err)
	} else {
		fmt.Printf("Take profit order placed successfully: %v\n", takeProfitOrder)
	}
}

// getOrderHistory retrieves and displays order history
func getOrderHistory(client *pi42.Client, symbol string) {
	// Get order history
	orders, err := client.Order.GetOrderHistory(pi42.OrderQueryParams{
		Symbol:   symbol,
		PageSize: 5, // Limit to 5 orders
	})
	if err != nil {
		log.Printf("Error getting order history: %v\n", err)
		return
	}

	fmt.Printf("Order History: Found %d orders\n", len(orders))
	for i, order := range orders {
		fmt.Printf("Order %d:\n", i+1)
		fmt.Printf("  ID: %v\n", order.ClientOrderID)
		fmt.Printf("  Symbol: %v\n", order.Symbol)
		fmt.Printf("  Type: %v\n", order.Type)
		fmt.Printf("  Side: %v\n", order.Side)
		fmt.Printf("  Price: %v\n", order.Price)
		fmt.Printf("  Quantity: %v\n", order.ExecutedQty)
		fmt.Printf("  Status: %v\n", order.Status)
		fmt.Printf("  Created At: %v\n", order.UpdatedAt)
	}

	// Get open orders
	openOrders, err := client.Order.GetOpenOrders(pi42.OrderQueryParams{
		Symbol: symbol,
	})
	if err != nil {
		log.Printf("Error getting open orders: %v\n", err)
		return
	}

	fmt.Printf("Open Orders: Found %d orders\n", len(openOrders))
	for i, order := range openOrders {
		fmt.Printf("Order %d:\n", i+1)
		fmt.Printf("  ID: %v\n", order.ClientOrderID)
		fmt.Printf("  Symbol: %v\n", order.Symbol)
		fmt.Printf("  Type: %v\n", order.Type)
		fmt.Printf("  Side: %v\n", order.Side)
		fmt.Printf("  Price: %v\n", order.Price)
		fmt.Printf("  Quantity: %v\n", order.OrderAmount)
	}
}

// cancelOrders demonstrates how to cancel orders
func cancelOrders(client *pi42.Client) {
	// Get open orders first
	openOrders, err := client.Order.GetOpenOrders(pi42.OrderQueryParams{})
	if err != nil {
		log.Printf("Error getting open orders: %v\n", err)
		return
	}

	if len(openOrders) == 0 {
		fmt.Println("No open orders to cancel")
		return
	}

	// Cancel the first open order
	firstOrder := openOrders[0]
	clientOrderID := firstOrder.ClientOrderID

	fmt.Printf("Cancelling order with client order ID: %s\n", clientOrderID)
	result, err := client.Order.DeleteOrder(clientOrderID)
	if err != nil {
		log.Printf("Error cancelling order: %v\n", err)
	} else {
		fmt.Printf("Order cancelled successfully: %v\n", result)
	}

	// Uncomment to test cancelling all orders
	/*
		fmt.Println("Cancelling all open orders...")
		result, err = client.Order.CancelAllOrders()
		if err != nil {
			log.Printf("Error cancelling all orders: %v\n", err)
		} else {
			fmt.Printf("All orders cancelled successfully: %v\n", result)
		}
	*/
}
