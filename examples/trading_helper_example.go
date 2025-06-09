package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/revanthstrakz/pi42"
)

func tradingHelperExample() {
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

	fmt.Println("=== Trading Helper Example ===")

	// Create a new TradingHelper for BTCINR with a 0.1% step size
	symbol := "BTCINR"
	helper, err := pi42.NewTradingHelper(client, symbol, 0.1)
	if err != nil {
		log.Fatalf("Error creating trading helper: %v", err)
	}

	// Display symbol information
	info := helper.GetSymbolInfo()
	fmt.Printf("\nSymbol Information for %s:\n", symbol)
	fmt.Printf("  Base Asset: %s\n", info.BaseAsset)
	fmt.Printf("  Quote Asset: %s\n", info.QuoteAsset)
	fmt.Printf("  Margin Asset: %s\n", info.MarginAsset)
	fmt.Printf("  Contract Type: %s\n", info.ContractType)
	fmt.Printf("  Min Quantity: %.8f %s\n", info.MinQuantity, info.BaseAsset)
	fmt.Printf("  Max Quantity: %.8f %s\n", info.MaxQuantity, info.BaseAsset)
	fmt.Printf("  Quantity Precision: %d\n", info.QuantityPrecision)
	fmt.Printf("  Price Precision: %d\n", info.PricePrecision)
	fmt.Printf("  Min Price Step: %.8f %s\n", info.MinPriceStep, info.QuoteAsset)

	// Get current market prices
	bestBid, bestAsk, err := helper.GetCurrentBestPrices()
	if err != nil {
		log.Fatalf("Error getting best prices: %v", err)
	}
	fmt.Printf("\nCurrent Market Prices:\n")
	fmt.Printf("  Best Bid: %.2f %s\n", bestBid, info.QuoteAsset)
	fmt.Printf("  Best Ask: %.2f %s\n", bestAsk, info.QuoteAsset)
	fmt.Printf("  Spread: %.2f %s (%.4f%%)\n",
		bestAsk-bestBid,
		info.QuoteAsset,
		(bestAsk-bestBid)*100/bestBid)

	// Calculate prices at various percentage differences
	fmt.Printf("\nPrice Levels:\n")

	// Buy prices (below best bid)
	fmt.Printf("  Buy Prices (below best bid):\n")
	for _, pct := range []float64{0.1, 0.5, 1.0, 2.0} {
		price, err := helper.CalculatePriceFromBestPrice(-pct)
		if err != nil {
			log.Printf("Error calculating price at %.1f%%: %v\n", -pct, err)
			continue
		}
		fmt.Printf("    %.1f%% below bid: %.2f %s\n", pct, price, info.QuoteAsset)
	}

	// Sell prices (above best ask)
	fmt.Printf("  Sell Prices (above best ask):\n")
	for _, pct := range []float64{0.1, 0.5, 1.0, 2.0} {
		price, err := helper.CalculatePriceFromBestPrice(pct)
		if err != nil {
			log.Printf("Error calculating price at +%.1f%%: %v\n", pct, err)
			continue
		}
		fmt.Printf("    %.1f%% above ask: %.2f %s\n", pct, price, info.QuoteAsset)
	}

	// Calculate order quantities for different investment amounts
	fmt.Printf("\nOrder Quantities for Different Investment Amounts:\n")
	for _, amount := range []float64{1000, 5000, 10000} {
		quantity, err := helper.CalculateOrderQuantity(amount)
		if err != nil {
			log.Printf("Error calculating quantity for %.2f %s: %v\n",
				amount, info.QuoteAsset, err)
			continue
		}
		fmt.Printf("  %.2f %s → %.8f %s\n",
			amount, info.QuoteAsset, quantity, info.BaseAsset)
	}

	fmt.Println("\n=== Trading Helper Example Completed ===")
}
