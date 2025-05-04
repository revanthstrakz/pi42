# Pi42 Go Client

## About

This Go client library provides a comprehensive and easy-to-use interface for the Pi42 cryptocurrency exchange API. It enables developers to build trading systems, market data applications, and automated trading bots with minimal effort. The library handles all the complexities of API authentication, rate limiting, WebSocket connections, and data processing, allowing you to focus on your trading strategies and application logic.

## Overview

The Pi42 Go Client provides a comprehensive interface to the Pi42 API, allowing developers to build trading applications, tools, and bots for the Pi42 exchange.

## Features

- Complete API coverage for Pi42 exchange
- Support for public and authenticated endpoints
- Real-time market data via Socketio
- Comprehensive error handling
- Easy-to-use API structure following Go conventions

## Installation

```bash
go get github.com/pi42/go-client
```

## Quick Start

```go
package main

import (
    "fmt"
    
    "github.com/pi42/go-client"
)

func main() {
    // Create a new client
    client := pi42.NewClient("YOUR_API_KEY", "YOUR_API_SECRET")
    
    // Get exchange info
    exchangeInfo, err := client.Exchange.ExchangeInfo("")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    contracts, _ := exchangeInfo["contracts"].([]interface{})
    fmt.Printf("Available contracts: %d\n", len(contracts))
    
    // Get ticker data
    ticker, err := client.Market.GetTicker24hr("BTCINR")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if data, ok := ticker["data"].(map[string]interface{}); ok {
        fmt.Printf("BTC Price: %v\n", data["c"])
    }
}
```

## API Structure

The client is organized into the following components:

- `Market`: Access to market data (tickers, klines, etc.)
- `Order`: Operations related to orders (place, cancel, etc.)
- `Position`: Operations related to positions (get, close, etc.)
- `Wallet`: Access to wallet information
- `Exchange`: Access to exchange information and settings
- `UserData`: Access to user-specific data
- `Socketio`: Real-time data streams

## Documentation

### Market API

```go
// Get 24-hour ticker data
ticker, err := client.Market.GetTicker24hr("BTCINR")

// Get candlestick (kline) data
klines, err := client.Market.GetKlines(pi42.KlinesParams{
    Pair:     "BTCINR",
    Interval: "1h",
    Limit:    5,
})

// Get aggregated trade data
trades, err := client.Market.GetAggTrades("BTCINR")

// Get order book depth
depth, err := client.Market.GetDepth("BTCINR")
```

### Order API

```go
// Place an order
order, err := client.Order.PlaceOrder(pi42.PlaceOrderParams{
    Symbol:      "BTCINR",
    Side:        "BUY",
    Type:        "LIMIT",
    Quantity:    0.01,
    Price:       4500000,
    MarginAsset: "INR",
})

// Get open orders
orders, err := client.Order.GetOpenOrders(pi42.OrderQueryParams{})

// Cancel an order
result, err := client.Order.DeleteOrder("CLIENT_ORDER_ID")
```

### Socketio API

```go
// Define a handler for ticker updates
handleTicker := func(data map[string]interface{}) {
    fmt.Printf("Ticker update: %v\n", data)
}

// Register the handler
client.Socketio.On("24hrTicker", handleTicker)

// Connect and subscribe to BTCINR ticker
client.Socketio.ConnectPublic([]string{"btcinr@ticker"})
```

## WebSocket Data Streams

The library provides a dedicated WebSocket client for subscribing to real-time market data streams.

### Creating a WebSocket Client

```go
import (
    "fmt"
    "time"
    
    "github.com/revanthstrakz/pi42"
    "github.com/zishang520/engine.io/v2/utils"
)

func main() {
    // Create a new WebSocket client
    client := pi42.NewSocketClient()
    
    // Start the client in a background goroutine
    go client.Init()
    
    // Use the client to subscribe to streams
    // ...
}
```

### Subscribing to Data Streams

```go
// Subscribe to specific topics
client.AddStream("btcinr@depth_0.1", "depthUpdate")
client.AddStream("btcinr@markPrice", "markPriceUpdate")
client.AddStream("btcinr@kline_1m", "kline")

// Wait for connection to establish
time.Sleep(2 * time.Second)

// Subscribe to more streams later as needed
client.AddStream("ethinr@markPrice", "markPriceUpdate")
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

### Unsubscribing from Streams

```go
// Remove a specific stream while keeping the event channel active
client.RemoveStream("btcinr@markPrice")
// The markPriceUpdate channel remains active for other subscribed topics
```

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

### Example Usage

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
    
    // Subscribe to topics
    client.AddStream("btcinr@depth_0.1", "depthUpdate")
    client.AddStream("btcinr@kline_1m", "kline")
    
    // Setup handlers for event channels
    klineChannel, exists := client.GetEventChannel("kline")
    if exists {
        go func() {
            for event := range klineChannel {
                fmt.Printf("Kline data: %v\n", event.Data)
            }
        }()
    }
    
    // Start the client
    go client.Init()
    
    // Keep the main process running
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

## Examples

Please refer to the examples directory for complete working examples.

## License

MIT License
