"""
Basic usage examples for the Pi42 API client.
"""

import os
import time
from pi42 import Pi42Client, Pi42APIError, Pi42RequestError

# Get API credentials from environment variables
API_KEY = os.environ.get("PI42_API_KEY")
API_SECRET = os.environ.get("PI42_API_SECRET")

# Create a client instance
client = Pi42Client(api_key=API_KEY, api_secret=API_SECRET)

def public_api_examples():
    """Examples of using the public API endpoints."""
    print("\n=== Public API Examples ===\n")
    
    # Get exchange info
    try:
        exchange_info = client.exchange.get_exchange_info()
        print(f"Exchange Info: Found {len(exchange_info.get('contracts', []))} contracts")
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error getting exchange info: {e}")
    
    # Get ticker data for BTC
    try:
        ticker = client.market.get_ticker_24hr("btc")
        print(f"BTC 24hr Ticker: Last price = {ticker.get('data', {}).get('c')}")
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error getting ticker: {e}")
    
    # Get klines data
    try:
        klines = client.market.get_klines("BTCINR", "1h", limit=5)
        print(f"BTCINR Klines: Retrieved {len(klines)} hourly candles")
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error getting klines: {e}")

def authenticated_api_examples():
    """Examples of using the authenticated API endpoints."""
    if not API_KEY or not API_SECRET:
        print("\n=== Authenticated API Examples (Skipped - No API Keys) ===\n")
        return
        
    print("\n=== Authenticated API Examples ===\n")
    
    # Get wallet details
    try:
        wallet = client.wallet.get_futures_wallet_details()
        print(f"Futures Wallet: Available balance = {wallet.get('withdrawableBalance')} INR")
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error getting wallet details: {e}")
    
    # Get open orders
    try:
        orders = client.order.get_open_orders()
        print(f"Open Orders: Found {len(orders)} open orders")
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error getting open orders: {e}")
    
    # Get open positions
    try:
        positions = client.position.get_positions("OPEN")
        print(f"Open Positions: Found {len(positions)} open positions")
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error getting positions: {e}")

def websocket_example():
    """Example of using the WebSocket API."""
    print("\n=== WebSocket Example ===\n")
    
    # Create a function to handle incoming data
    def handle_ticker(data):
        print(f"Received ticker update for {data.get('s')}: Price = {data.get('c')}")
    
    # Using the public WebSocket API
    ws_client = client.websocket
    ws_client.on('24hrTicker', handle_ticker)
    
    print("Connecting to WebSocket and subscribing to BTCINR ticker...")
    ws_client.connect_public(topics=["btcinr@ticker"])
    
    # Wait for some data
    print("Waiting for data (press Ctrl+C to exit)...")
    try:
        time.sleep(30)  # Wait for 30 seconds to receive some data
    except KeyboardInterrupt:
        pass
    finally:
        ws_client.close()
        print("WebSocket connection closed")

if __name__ == "__main__":
    public_api_examples()
    authenticated_api_examples()
    websocket_example()
