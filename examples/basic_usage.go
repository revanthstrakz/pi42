package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/revanthstrakz/pi42"

	"github.com/joho/godotenv"
)

func main() {
	// BasicUsageExample()
	// DataStreamExample()
	LimitCheck()
}

func BasicUsageExample() {
	// Load environment variables from .env file
	err := godotenv.Load(filepath.Join(".", ".env"))
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
	// OrderExample(client)
}

func publicAPIExamples(client *pi42.Client) {
	fmt.Println("\n=== Public API Examples ===")

	// Get exchange info with structured response
	exchangeInfo, err := client.Exchange.ExchangeInfo("")
	if err != nil {
		fmt.Printf("Error getting exchange info: %v\n", err)
	} else {
		fmt.Printf("Exchange Info: Found %d contracts\n", len(exchangeInfo.Contracts))
		fmt.Printf("Available markets: %v\n", exchangeInfo.Markets)
		fmt.Printf("Asset precisions: %v\n", exchangeInfo.AssetPrecisions)

		// Print detailed information about the first contract
		if len(exchangeInfo.Contracts) > 0 {
			contract := exchangeInfo.Contracts[0]
			fmt.Printf("\nDetails for %s (%s):\n", contract.Name, contract.ContractName)
			fmt.Printf("  Base/Quote: %s/%s\n", contract.BaseAsset, contract.QuoteAsset)
			fmt.Printf("  Price Precision: %s\n", contract.PricePrecision)
			fmt.Printf("  Quantity Precision: %s\n", contract.QuantityPrecision)
			fmt.Printf("  Max Leverage: %s\n", contract.MaxLeverage)
			fmt.Printf("  Supported Margin Assets: %v\n", contract.MarginAssetsSupported)
			fmt.Printf("  Order Types: %v\n", contract.OrderTypes)

			if len(contract.Filters) > 0 {
				fmt.Printf("  Filters:\n")
				for _, filter := range contract.Filters {
					fmt.Printf("    %s: ", filter.FilterType)
					if filter.MinQty != "" {
						fmt.Printf("Min=%s, ", filter.MinQty)
					}
					if filter.MaxQty != "" {
						fmt.Printf("Max=%s, ", filter.MaxQty)
					}
					if filter.Notional != "" {
						fmt.Printf("Notional=%s, ", filter.Notional)
					}
					if filter.Limit != "" {
						fmt.Printf("Limit=%s", filter.Limit)
					}
					fmt.Println()
				}
			}
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
			fmt.Printf("First candle: Open=%v, Close=%v\n", klines[0].Open, klines[0].Close)
		}
	}

	// Example of using the depth API with structured response
	depthData, err := client.Market.GetDepth("BTCINR")
	if err != nil {
		fmt.Printf("Error getting order book depth: %v\n", err)
	} else {
		fmt.Printf("\nBTCINR Order Book: %d bids, %d asks\n",
			len(depthData.Data.Bids), len(depthData.Data.Asks))

		// Show the top bid and ask
		if len(depthData.Data.Bids) > 0 {
			fmt.Printf("Best bid: %s @ %s\n",
				depthData.Data.Bids[0][0], depthData.Data.Bids[0][1])
		}

		if len(depthData.Data.Asks) > 0 {
			fmt.Printf("Best ask: %s @ %s\n",
				depthData.Data.Asks[0][0], depthData.Data.Asks[0][1])
		}
	}
}

func authenticatedAPIExamples(client *pi42.Client) {
	fmt.Println("\n=== Authenticated API Examples ===")

	if client.APIKey == "" || client.APISecret == "" {
		fmt.Println("Skipped - No API Keys")
		return
	}

	// Get wallet details with structured response
	futuresWallet, err := client.Wallet.FuturesWalletDetails("INR")
	if err != nil {
		fmt.Printf("Error getting futures wallet details: %v\n", err)
	} else {
		fmt.Printf("Futures Wallet: Available balance = %s INR\n", futuresWallet.WithdrawableBalance)
		fmt.Printf("  Total Balance: %s INR\n", futuresWallet.WalletBalance)
		fmt.Printf("  Locked Balance: %s INR\n", futuresWallet.LockedBalance)
		fmt.Printf("  Margin Balance: %s INR\n", futuresWallet.MarginBalance)
		fmt.Printf("  Unrealized PnL (Cross): %s INR\n", futuresWallet.UnrealisedPnlCross)
		fmt.Printf("  Unrealized PnL (Isolated): %s INR\n", futuresWallet.UnrealisedPnlIsolated)
	}

	fundingWallet, err := client.Wallet.FundingWalletDetails("INR")
	if err != nil {
		fmt.Printf("Error getting funding wallet details: %v\n", err)
	} else {
		fmt.Printf("Funding Wallet: Available balance = %s INR\n", fundingWallet.WithdrawableBalance)
		fmt.Printf("  Total Balance: %s INR\n", fundingWallet.WalletBalance)
		fmt.Printf("  Locked Balance: %s INR\n", fundingWallet.LockedBalance)
	}

	// Try updating leverage for a contract
	leverageUpdate, err := client.Exchange.UpdateLeverage(10, "BTCINR")
	if err != nil {
		fmt.Printf("Error updating leverage: %v\n", err)
	} else {
		fmt.Printf("Leverage updated for %s: %d\n",
			leverageUpdate.ContractName,
			leverageUpdate.UpdatedLeverage)
	}

	// Get open orders using structured response
	openOrders, err := client.Order.GetOpenOrders(pi42.OrderQueryParams{})
	if err != nil {
		fmt.Printf("Error getting open orders: %v\n", err)
	} else {
		fmt.Printf("Open Orders: Found %d open orders\n", len(openOrders))

		// Display some details about open orders if available
		if len(openOrders) > 0 {
			for i, order := range openOrders {
				fmt.Printf("  Order %d:\n", i+1)
				fmt.Printf("    Symbol: %s\n", order.Symbol)
				fmt.Printf("    Type: %s\n", order.Type)
				fmt.Printf("    Side: %s\n", order.Side)
				fmt.Printf("    Price: %.2f\n", order.Price)
				fmt.Printf("    Amount: %.8f\n", order.OrderAmount)
				fmt.Printf("    Status: %s\n", order.Status)
			}
		}
	}

	// Get trade history using structured response
	trades, err := client.UserData.GetTradeHistory(pi42.DataQueryParams{
		PageSize: 5,
	})
	if err != nil {
		fmt.Printf("Error getting trade history: %v\n", err)
	} else {
		fmt.Printf("Trade History: Found %d trades\n", len(trades))

		// Display some details about trades if available
		if len(trades) > 0 {
			for i, trade := range trades {
				fmt.Printf("  Trade %d: %s %s %s at %.8f\n",
					i+1, trade.Side, trade.Type, trade.Symbol, trade.Price)
			}
		}
	}

	// Get open positions using structured response
	positions, err := client.Position.GetPositions("OPEN", pi42.PositionQueryParams{})
	if err != nil {
		fmt.Printf("Error getting positions: %v\n", err)
	} else {
		fmt.Printf("Open Positions: Found %d open positions\n", len(positions))

		// Display some details about positions if available
		for i, pos := range positions {
			fmt.Printf("  Position %d: %s %s\n", i+1, pos.PositionType, pos.ContractPair)
			fmt.Printf("    Entry Price: %.2f\n", pos.EntryPrice)
			fmt.Printf("    Size: %.8f\n", pos.PositionAmount)
			fmt.Printf("    Margin: %.2f %s\n", pos.Margin, pos.MarginAsset)
			fmt.Printf("    Leverage: %dx\n", pos.Leverage)
		}
	}
}
