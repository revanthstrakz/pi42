package main

import (
	"fmt"
	"log"

	"github.com/revanthstrakz/pi42"
)

func OrderExample(client *pi42.Client) {
	orderAPI := client.Order

	// Example 1: Using Bullet function for market order with structured response
	fmt.Println("\nExample 1: Using Bullet function for market order")
	marketOrder, err := orderAPI.Bullet(pi42.BulletParams{
		Symbol:     "ALCHUSDT",
		Side:       "BUY",
		OrderType:  "MARKET",
		Count:      1.0,
		ReduceOnly: false,
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		// Access fields directly from the structured response
		fmt.Printf("Order placed successfully!\n")
		fmt.Printf("  Order ID: %v\n", marketOrder.ID)
		fmt.Printf("  Client Order ID: %s\n", marketOrder.ClientOrderID)
		fmt.Printf("  Symbol: %s\n", marketOrder.Symbol)
		fmt.Printf("  Type: %s\n", marketOrder.Type)
		fmt.Printf("  Side: %s\n", marketOrder.Side)
		fmt.Printf("  Quantity: %f\n", marketOrder.OrderAmount)
		fmt.Printf("  Available Balance: %f %s\n", marketOrder.AvailableBalance, marketOrder.MarginAsset)
		fmt.Printf("  Leverage: %dx\n", marketOrder.Leverage)
	}

	//  get position id for ALCHUSDT
	positions, err := client.Position.GetPositions("OPEN", pi42.PositionQueryParams{})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	positionID := ""
	for _, position := range positions {
		if position.ContractPair == "ALCHUSDT" {
			positionID = position.PositionID
			break
		}
	}

	// Example 2: Using Bullet function for limit order
	fmt.Println("\nExample 2: Using Bullet function for limit order")
	limitOrder, err := orderAPI.Bullet(pi42.BulletParams{
		Symbol:     "ALCHUSDT",
		Side:       "BUY",
		OrderType:  "LIMIT",
		Price:      0.200,
		Count:      2.5,
		ReduceOnly: true,
		PositionID: positionID,
	})
	printStructuredResponse("Limit Order", limitOrder, err)

	// Example 3: Using Bullet function for stop market order
	fmt.Println("\nExample 3: Using Bullet function for stop market order")
	stopMarketOrder, err := orderAPI.Bullet(pi42.BulletParams{
		Symbol:     "ALCHUSDT",
		Side:       "SELL",
		OrderType:  "STOP_MARKET",
		StopPrice:  11.0, // Trigger price
		Count:      1.0,
		ReduceOnly: true,
	})
	printStructuredResponse("Stop Market Order", stopMarketOrder, err)

	// Example 4: Using Bullet function for stop limit order
	fmt.Println("\nExample 4: Using Bullet function for stop limit order")
	stopLimitOrder, err := orderAPI.Bullet(pi42.BulletParams{
		Symbol:     "ALCHUSDT",
		Side:       "SELL",
		OrderType:  "STOP_LIMIT",
		Price:      10.8, // Execution price
		StopPrice:  11.0, // Trigger price
		Count:      1.5,
		ReduceOnly: true,
	})
	printStructuredResponse("Stop Limit Order", stopLimitOrder, err)

}

// Helper function to print structured responses
func printStructuredResponse(orderType string, response *pi42.OrderResponse, err error) {
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%s placed successfully!\n", orderType)
	fmt.Printf("  Order ID: %v\n", response.ID)
	fmt.Printf("  Client Order ID: %s\n", response.ClientOrderID)
	fmt.Printf("  Type: %s\n", response.Type)
	fmt.Printf("  Side: %s\n", response.Side)
	fmt.Printf("  Price: %f\n", response.Price)

	// For stop orders, show the stop price field
	if response.Type == "STOP_MARKET" || response.Type == "STOP_LIMIT" {
		fmt.Printf("  Stop Price: %f\n", response.StopPrice)
	}

	fmt.Printf("  Quantity: %f\n", response.OrderAmount)
	fmt.Printf("  Locked Margin: %f %s\n", response.LockedMargin, response.MarginAsset)
	fmt.Printf("  Time: %s\n", response.Time)
}
