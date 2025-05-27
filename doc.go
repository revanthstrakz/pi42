// Package pi42 provides a comprehensive Go client for the Pi42 cryptocurrency exchange API.
//
// This package enables developers to:
//   - Access market data (tickers, candlesticks, order books)
//   - Place and manage trading orders
//   - Monitor and manage positions
//   - Access wallet information
//   - Stream real-time data via WebSockets
//   - And more
//
// The client is designed to be easy to use while providing all the functionality
// needed to build trading systems, market data applications, and automated trading bots.
//
// Basic usage:
//
//	import "github.com/revanthstrakz/pi42"
//
//	// Create client
//	client := pi42.NewClient("API_KEY", "API_SECRET")
//
//	// Access market data
//	ticker, err := client.Market.GetTicker24hr("BTCINR")
//
//	// Place orders
//	order, err := client.Order.Bullet(pi42.BulletParams{
//		Symbol:    "BTCINR",
//		Side:      "BUY",
//		OrderType: "MARKET",
//		Count:     1.0,
//	})
//
// See the README.md file for more detailed documentation.
package pi42
