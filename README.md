# Pi42 Go Client

## About

This Go client library provides a comprehensive and easy-to-use interface for the Pi42 cryptocurrency exchange API. It enables developers to build trading systems, market data applications, and automated trading bots with minimal effort. The library handles all the complexities of API authentication, rate limiting, WebSocket connections, and data processing, allowing you to focus on your trading strategies and application logic.

## Installation

```bash
go get github.com/revanthstrakz/pi42
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/revanthstrakz/pi42"
)

func main() {
    // Get API credentials from environment variables
    apiKey := os.Getenv("PI42_API_KEY")
    apiSecret := os.Getenv("PI42_API_SECRET")
    
    // Create a new client
    client := pi42.NewClient(apiKey, apiSecret)
    
    // Get exchange info
    exchangeInfo, err := client.Exchange.ExchangeInfo("")
    if err != nil {
        log.Fatalf("Error: %v\n", err)
    }
    
    contracts, _ := exchangeInfo["contracts"].([]interface{})
    fmt.Printf("Available contracts: %d\n", len(contracts))
    
    // Get ticker data
    ticker, err := client.Market.GetTicker24hr("BTCINR")
    if err != nil {
        log.Fatalf("Error: %v\n", err)
    }
    
    if data, ok := ticker["data"].(map[string]interface{}); ok {
        fmt.Printf("BTC Price: %v\n", data["c"])
    }
}
```

## Client Structure

The client is organized into the following components:

- `Market`: Access to market data (tickers, klines, etc.)
- `Order`: Operations related to orders (place, cancel, etc.)
- `Position`: Operations related to positions (get, close, etc.)
- `Wallet`: Access to wallet information
- `Exchange`: Access to exchange information and settings
- `UserData`: Access to user-specific data

## Authentication

To use authenticated endpoints, you need to provide your API key and secret:

```go
client := pi42.NewClient("YOUR_API_KEY", "YOUR_API_SECRET")
```

For security reasons, it's recommended to use environment variables:

```go
apiKey := os.Getenv("PI42_API_KEY")
apiSecret := os.Getenv("PI42_API_SECRET")
client := pi42.NewClient(apiKey, apiSecret)
```

You can also use a .env file with the help of godotenv:

```go
import "github.com/joho/godotenv"

func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file")
    }
}
```

## Detailed API Usage

### Market API

The Market API provides access to market data such as tickers, candlesticks, trades, and order book data.

```go
// Get 24-hour ticker data
ticker, err := client.Market.GetTicker24hr("BTCINR")

// Get candlestick (kline) data
klines, err := client.Market.GetKlines(pi42.KlinesParams{
    Pair:     "BTCINR",
    Interval: "1h",
    Limit:    10,
    // Optional parameters
    // StartTime: 1625000000000,
    // EndTime:   1625100000000,
})

// Get aggregated trade data
trades, err := client.Market.GetAggTrades("BTCINR")

// Get order book depth
depth, err := client.Market.GetDepth("BTCINR")
```

### Order API

The Order API allows you to place, query, and cancel orders.

#### Simple Order Placement with Bullet

The `Bullet` method provides a simplified way to place orders with automatic handling of precision and other parameters:

```go
// Market buy order
marketOrder, err := client.Order.Bullet(pi42.BulletParams{
    Symbol:     "BTCINR",
    Side:       "BUY",
    OrderType:  "MARKET",
    Count:      1.0,     // This is a multiplier of the minimum quantity
    ReduceOnly: false,
})

// Limit buy order
limitOrder, err := client.Order.Bullet(pi42.BulletParams{
    Symbol:     "BTCINR",
    Side:       "BUY",
    OrderType:  "LIMIT",
    Price:      4500000, // Limit price
    Count:      2.5,     // This is a multiplier of the minimum quantity
    ReduceOnly: false,
})

// Stop market sell order
stopMarketOrder, err := client.Order.Bullet(pi42.BulletParams{
    Symbol:     "BTCINR",
    Side:       "SELL",
    OrderType:  "STOP_MARKET",
    StopPrice:  4400000, // Trigger price
    Count:      1.0,
    ReduceOnly: true,
})

// Stop limit sell order
stopLimitOrder, err := client.Order.Bullet(pi42.BulletParams{
    Symbol:     "BTCINR",
    Side:       "SELL",
    OrderType:  "STOP_LIMIT",
    Price:      4390000, // Execution price
    StopPrice:  4400000, // Trigger price
    Count:      1.5,
    ReduceOnly: true,
})
```

#### Advanced Order Placement

For more control, you can use the `PlaceOrder` method directly:

```go
order, err := client.Order.PlaceOrder(pi42.PlaceOrderParams{
    Symbol:          "BTCINR",
    Side:            "BUY",
    Type:            "LIMIT",
    Quantity:        0.01,
    Price:           4500000,
    MarginAsset:     "INR",
    ReduceOnly:      false,
    TakeProfitPrice: 4600000, // Optional
    StopLossPrice:   4400000, // Optional
})
```

#### Query Orders

```go
// Get open orders
openOrders, err := client.Order.GetOpenOrders(pi42.OrderQueryParams{
    Symbol:    "BTCINR", // Optional - omit for all symbols
    PageSize:  50,       // Optional
    SortOrder: "DESC",   // Optional
})

// Get order history
orderHistory, err := client.Order.GetOrderHistory(pi42.OrderQueryParams{
    Symbol:         "BTCINR", // Optional
    StartTimestamp: 1625000000000, // Optional
    EndTimestamp:   1625100000000, // Optional
    PageSize:       100, // Optional
})
```

#### Cancel Orders

```go
// Cancel a specific order
result, err := client.Order.DeleteOrder("CLIENT_ORDER_ID")

// Cancel all open orders
result, err := client.Order.CancelAllOrders()
```

### Position API

The Position API allows you to manage your trading positions.

```go
// Get open positions
openPositions, err := client.Position.GetPositions("OPEN", pi42.PositionQueryParams{
    Symbol: "BTCINR", // Optional - omit for all symbols
})

// Get closed positions
closedPositions, err := client.Position.GetPositions("CLOSED", pi42.PositionQueryParams{
    StartTimestamp: 1625000000000, // Optional
    EndTimestamp:   1625100000000, // Optional
})

// Get details of a specific position
positionDetails, err := client.Position.GetPosition("POSITION_ID")

// Close all positions
result, err := client.Position.CloseAllPositions()
```

### Wallet API

The Wallet API allows you to access your wallet information.

```go
// Get futures wallet details
futuresWallet, err := client.Wallet.FuturesWalletDetails("INR")

// Get funding wallet details
fundingWallet, err := client.Wallet.FundingWalletDetails("INR")
```

### Exchange API

The Exchange API provides information about the exchange and lets you update certain trading preferences.

```go
// Get exchange info for all markets
exchangeInfo, err := client.Exchange.ExchangeInfo("")

// Get exchange info for a specific market
exchangeInfo, err := client.Exchange.ExchangeInfo("futures")

// Update leverage for a contract
result, err := client.Exchange.UpdateLeverage(10, "BTCINR")

// Update both leverage and margin mode (ISOLATED or CROSS)
result, err := client.Exchange.UpdatePreference(10, "ISOLATED", "BTCINR")
```

### User Data API

The User Data API provides access to user-specific data.

```go
// Get trade history
tradeHistory, err := client.UserData.GetTradeHistory(pi42.DataQueryParams{
    Symbol:         "BTCINR", // Optional
    StartTimestamp: 1625000000000, // Optional
    EndTimestamp:   1625100000000, // Optional
    PageSize:       100, // Optional
})

// Get transaction history
txHistory, err := client.UserData.GetTransactionHistory(pi42.TransactionHistoryParams{
    DataQueryParams: pi42.DataQueryParams{
        Symbol:    "BTCINR", // Optional
        PageSize:  50,       // Optional
    },
    PositionID: "POSITION_ID", // Optional
})

// Get listen key for WebSocket user data stream
listenKey, err := client.UserData.CreateListenKey()
```

## WebSocket Data Streams

The library provides a dedicated WebSocket client for subscribing to real-time market data streams.

### Creating a WebSocket Client

```go
// Create a new WebSocket client
client := pi42.NewSocketClient()

// Start the client in a background goroutine
go client.Init()
```

### Subscribing to Data Streams

```go
// Subscribe to specific topics
client.AddStream("btcinr@depth_0.1", "depthUpdate")
client.AddStream("btcinr@markPrice", "markPriceUpdate")
client.AddStream("btcinr@kline_1m", "kline")
```

### Receiving Data via Channels

```go
// Get channel for a specific event type
markPriceChannel, exists := client.GetEventChannel("markPriceUpdate")
if exists {
    go func() {
        for event := range markPriceChannel {
            fmt.Printf("Mark price update received: %v\n", event.Data)
        }
    }()
}

// Get channel for another event type
klineChannel, exists := client.GetEventChannel("kline")
if exists {
    go func() {
        for event := range klineChannel {
            fmt.Printf("Kline data received - Topic: %s, Data: %v\n", 
                event.Topic, event.Data)
        }
    }()
}
```

### Supported WebSocket Topics

The format for topics is: `<symbol>@<channel>_<options>`

Examples:
- `btcinr@depth_0.1` - Order book updates for BTCINR with 0.1 granularity
- `btcinr@markPrice` - Mark price updates for BTCINR
- `btcinr@kline_1m` - 1-minute candlestick data for BTCINR (intervals: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M)
- `btcinr@trade` - Trade data for BTCINR
- `btcinr@ticker` - Ticker data for BTCINR

### Supported Event Types

The WebSocket client supports the following event types:

- `depthUpdate`: Order book updates
- `markPriceUpdate`: Mark price updates
- `kline`: Candlestick/kline data
- `aggTrade`: Aggregated trade data
- `24hrTicker`: 24-hour ticker updates
- `marketInfo`: Market information
- `tickerArr`: Array of ticker data
- `markPriceArr`: Array of mark prices
- `allContractDetails`: Contract detail updates

## Complete Examples

### Basic Market Data Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/revanthstrakz/pi42"
)

func main() {
    // Create a client without authentication for public endpoints
    client := pi42.NewClient("", "")
    
    // Get ticker data for BTC
    ticker, err := client.Market.GetTicker24hr("BTCINR")
    if err != nil {
        log.Fatalf("Error getting ticker: %v", err)
    }
    
    data, ok := ticker["data"].(map[string]interface{})
    if ok {
        fmt.Printf("BTC Last Price: %v\n", data["c"])
        fmt.Printf("BTC 24h High: %v\n", data["h"])
        fmt.Printf("BTC 24h Low: %v\n", data["l"])
        fmt.Printf("BTC 24h Volume: %v\n", data["v"])
    }
    
    // Get recent candlestick data
    klines, err := client.Market.GetKlines(pi42.KlinesParams{
        Pair:     "BTCINR",
        Interval: "1h",
        Limit:    5,
    })
    if err != nil {
        log.Fatalf("Error getting klines: %v", err)
    }
    
    fmt.Println("\nRecent 1h Candles:")
    for i, candle := range klines {
        fmt.Printf("Candle %d - Open: %v, High: %v, Low: %v, Close: %v\n",
            i+1, candle["open"], candle["high"], candle["low"], candle["close"])
    }
}
```

### Trading Example

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/joho/godotenv"
    "github.com/revanthstrakz/pi42"
)

func init() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: Error loading .env file")
    }
}

func main() {
    // Get API credentials from environment variables
    apiKey := os.Getenv("PI42_API_KEY")
    apiSecret := os.Getenv("PI42_API_SECRET")
    
    // Create a client with authentication
    client := pi42.NewClient(apiKey, apiSecret)
    
    // Check wallet balance
    wallet, err := client.Wallet.FuturesWalletDetails("INR")
    if err != nil {
        log.Fatalf("Error getting wallet details: %v", err)
    }
    
    fmt.Printf("Available balance: %v INR\n", wallet["withdrawableBalance"])
    
    // Place a market order
    marketOrder, err := client.Order.Bullet(pi42.BulletParams{
        Symbol:     "BTCINR",
        Side:       "BUY",
        OrderType:  "MARKET",
        Count:      1.0,
        ReduceOnly: false,
    })
    
    if err != nil {
        log.Fatalf("Error placing market order: %v", err)
    }
    
    fmt.Printf("Market order placed successfully!\n")
    fmt.Printf("Order ID: %v\n", marketOrder.ID)
    fmt.Printf("Symbol: %s\n", marketOrder.Symbol)
    fmt.Printf("Quantity: %f\n", marketOrder.OrderAmount)
    
    // Get open positions
    positions, err := client.Position.GetPositions("OPEN", pi42.PositionQueryParams{})
    if err != nil {
        log.Fatalf("Error getting positions: %v", err)
    }
    
    fmt.Printf("\nOpen Positions: %d\n", len(positions))
    for _, pos := range positions {
        fmt.Printf("  Symbol: %s, Size: %f, Entry Price: %f, Margin: %f %s\n",
            pos.ContractPair, pos.PositionSize, pos.EntryPrice, pos.Margin, pos.MarginAsset)
    }
}
```

### WebSocket Real-time Data Example

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/revanthstrakz/pi42"
)

func main() {
    // Create a new WebSocket client
    client := pi42.NewSocketClient()
    
    // Add streams before starting the client
    client.AddStream("btcinr@depth_0.1", "depthUpdate")
    client.AddStream("btcinr@markPrice", "markPriceUpdate")
    client.AddStream("btcinr@kline_1m", "kline")
    
    // Get channels for different data types
    markPriceChannel, markPriceExists := client.GetEventChannel("markPriceUpdate")
    klineChannel, klineExists := client.GetEventChannel("kline")
    
    // Setup handlers for event channels
    if markPriceExists {
        go func() {
            for event := range markPriceChannel {
                fmt.Printf("Mark price update: %v\n", event.Data)
            }
        }()
    }
    
    if klineExists {
        go func() {
            for event := range klineChannel {
                fmt.Printf("Kline data - Topic: %s, Data: %v\n", 
                    event.Topic, event.Data)
            }
        }()
    }
    
    // Start the client in a separate goroutine
    go client.Init()
    
    // Add more streams after a delay
    time.Sleep(2 * time.Second)
    client.AddStream("ethinr@markPrice", "markPriceUpdate")
    
    // The program will continue to receive data via WebSocket
    // Keep the main goroutine running
    select {}
}
```

## Error Handling

The API uses Go's error handling approach and provides two types of errors:

- `APIError`: For errors returned by the Pi42 API
- Standard Go errors: For network and other client-side errors

Example:

```go
ticker, err := client.Market.GetTicker24hr("BTCINR")
if err != nil {
    if apiErr, ok := err.(pi42.APIError); ok {
        fmt.Printf("API Error Code: %d, Message: %s\n", apiErr.ErrorCode, apiErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Best Practices

1. **Rate Limiting**: Be mindful of API rate limits, especially for authenticated endpoints.
2. **Error Handling**: Always check errors and handle them appropriately.
3. **Secure Credentials**: Store API keys and secrets securely, and never hardcode them in your application.
4. **Reconnection Logic**: For WebSocket connections, implement reconnection logic to handle disconnections.
5. **Graceful Shutdown**: Implement proper shutdown procedures for WebSocket connections.

## License

MIT License
