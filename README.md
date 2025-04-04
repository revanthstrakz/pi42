# Pi42 API Client

A Python client for interacting with the Pi42 API. This client provides easy access to both public and authenticated endpoints.

## Installation

```bash
pip install pi42-api
```

Or install from the source:

```bash
git clone https://github.com/username/pi42-api.git
cd pi42-api
pip install .
```

## Features

- Access to all Pi42 API endpoints (market data, order management, wallet, etc.)
- Support for WebSocket connections for real-time data
- Error handling and type annotations
- Well-documented API methods
- Example scripts to demonstrate usage

## Getting Started

### Public API Access

```python
from pi42 import Pi42Client

# Create a client instance for public endpoints (no API key required)
client = Pi42Client()

# Get 24-hour ticker data for BTC
ticker = client.market.get_ticker_24hr("btc")
print(f"BTC price: {ticker.get('data', {}).get('c')}")

# Get klines (candlestick) data
klines = client.market.get_klines("BTCINR", "1h", limit=10)
print(f"Retrieved {len(klines)} hourly candles")
```

### Authenticated API Access

```python
from pi42 import Pi42Client

# Create a client instance with API key and secret
client = Pi42Client(api_key="your_api_key", api_secret="your_api_secret")

# Get wallet details
wallet = client.wallet.get_futures_wallet_details()
print(f"Available balance: {wallet.get('withdrawableBalance')} INR")

# Place an order
order = client.order.place_order(
    symbol="BTCINR",
    side="BUY",
    order_type="LIMIT",
    quantity=0.001,
    price=5000000,
)
print(f"Order placed: {order.get('clientOrderId')}")
```

### WebSocket Example

```python
import time
from pi42 import Pi42Client

# Create a client instance
client = Pi42Client()

# Define a callback function for ticker updates
def handle_ticker(data):
    print(f"Ticker update for {data.get('s')}: Price = {data.get('c')}")

# Set up the WebSocket client
ws_client = client.websocket
ws_client.on('24hrTicker', handle_ticker)

# Connect and subscribe to BTCINR ticker
ws_client.connect_public(topics=["btcinr@ticker"])

# Keep the script running to receive updates
try:
    time.sleep(60)  # Run for 60 seconds
finally:
    ws_client.close()
```

## API Modules

- `client`: Main API client with request handling and authentication
- `market`: Market data endpoints (klines, tickers, depth)
- `order`: Order management endpoints
- `position`: Position management endpoints
- `wallet`: Wallet and balance endpoints
- `exchange`: Exchange settings and configuration
- `user_data`: User data endpoints
- `websocket`: WebSocket connections for real-time data

## Error Handling

The client provides two main exception types:

- `Pi42APIError`: Raised when the API returns an error response
- `Pi42RequestError`: Raised for network issues or request problems

```python
from pi42 import Pi42Client, Pi42APIError, Pi42RequestError

client = Pi42Client()

try:
    ticker = client.market.get_ticker_24hr("invalid_symbol")
except Pi42APIError as e:
    print(f"API Error: {e}")
except Pi42RequestError as e:
    print(f"Request Error: {e}")
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
