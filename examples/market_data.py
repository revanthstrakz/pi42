"""
Examples of accessing market data with the Pi42 API client.
"""

import json
from datetime import datetime
import pandas as pd
from pi42 import Pi42Client, Pi42APIError, Pi42RequestError

# Create a client instance (no auth needed for public endpoints)
client = Pi42Client()

def format_timestamp(timestamp):
    """Convert a millisecond timestamp to a human-readable date."""
    return datetime.fromtimestamp(int(timestamp) / 1000).strftime('%Y-%m-%d %H:%M:%S')

def get_ticker_data():
    """Get and display 24-hour ticker data for multiple cryptocurrencies."""
    print("\n=== 24-Hour Ticker Data ===\n")
    
    cryptocurrencies = ["btc", "eth", "sol", "xrp"]
    
    for crypto in cryptocurrencies:
        try:
            data = client.market.get_ticker_24hr(crypto)
            ticker = data.get('data', {})
            
            print(f"Symbol: {ticker.get('s')}")
            print(f"Last Price: {ticker.get('c')}")
            print(f"24h Change: {ticker.get('p')} ({ticker.get('P')}%)")
            print(f"24h High: {ticker.get('h')}")
            print(f"24h Low: {ticker.get('l')}")
            print(f"24h Volume: {ticker.get('v')}")
            print("-" * 50)
            
        except (Pi42APIError, Pi42RequestError) as e:
            print(f"Error fetching {crypto} ticker: {e}")

def get_order_book():
    """Get and display order book data."""
    print("\n=== Order Book Data ===\n")
    
    try:
        data = client.market.get_depth("btc")
        depth = data.get('data', {})
        
        print(f"Symbol: {depth.get('s')}")
        print(f"Last Update ID: {depth.get('u')}")
        
        print("\nTop 5 Bids:")
        bids = depth.get('b', [])[:5]
        for i, bid in enumerate(bids, 1):
            print(f"{i}. Price: {bid[0]}, Quantity: {bid[1]}")
        
        print("\nTop 5 Asks:")
        asks = depth.get('a', [])[:5]
        for i, ask in enumerate(asks, 1):
            print(f"{i}. Price: {ask[0]}, Quantity: {ask[1]}")
            
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error fetching order book: {e}")

def get_klines_as_dataframe():
    """Get kline data and convert it to a pandas DataFrame."""
    print("\n=== Kline Data as DataFrame ===\n")
    
    try:
        klines = client.market.get_klines("BTCINR", "1h", limit=24)
        
        # Convert klines data to DataFrame
        df = pd.DataFrame(klines, columns=[
            "startTime", "open", "high", "low", "close", "endTime", "volume"
        ])
        
        # Convert timestamps to datetime
        df["startTime"] = pd.to_datetime(df["startTime"].astype(float), unit='ms')
        df["endTime"] = pd.to_datetime(df["endTime"].astype(float), unit='ms')
        
        # Convert numeric columns to float
        for col in ["open", "high", "low", "close", "volume"]:
            df[col] = df[col].astype(float)
        
        print(df.head())
        
        # Example of calculating simple moving average
        df["SMA_5"] = df["close"].rolling(window=5).mean()
        
        print("\nWith 5-period Simple Moving Average:")
        print(df[["startTime", "close", "SMA_5"]].tail(10))
        
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error fetching klines: {e}")
    except ImportError:
        print("pandas is not installed. Install with 'pip install pandas'")

def get_recent_trades():
    """Get and display recent aggregated trades."""
    print("\n=== Recent Aggregated Trades ===\n")
    
    try:
        data = client.market.get_agg_trades("btc")
        trades = data.get('data', [])
        
        if not trades:
            print("No recent trades found")
            return
            
        print(f"Most recent trade for {trades[0].get('s')}:")
        trade = trades[0]
        print(f"Price: {trade.get('p')}")
        print(f"Quantity: {trade.get('q')}")
        print(f"Time: {format_timestamp(trade.get('T'))}")
        print(f"Market Maker: {'Yes' if trade.get('m') else 'No'}")
        
        print("\nAll recent trades:")
        for i, trade in enumerate(trades, 1):
            print(f"{i}. Price: {trade.get('p')}, Quantity: {trade.get('q')}, Time: {format_timestamp(trade.get('T'))}")
            
    except (Pi42APIError, Pi42RequestError) as e:
        print(f"Error fetching recent trades: {e}")

if __name__ == "__main__":
    get_ticker_data()
    get_order_book()
    get_klines_as_dataframe()
    get_recent_trades()
