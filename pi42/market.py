"""
Market data API endpoints for Pi42.
"""

from typing import Dict, Any, Optional, List, Union


class MarketAPI:
    """
    Market API client for accessing Pi42 market data endpoints.
    """

    def __init__(self, client):
        """
        Initialize the Market API client.

        Args:
            client: The main Pi42Client instance
        """
        self._client = client

    def get_ticker_24hr(self, contract_pair: str) -> Dict[str, Any]:
        """
        Get 24-hour ticker data for a specific trading pair.

        Args:
            contract_pair: Trading pair symbol (e.g., 'btc' for BTCINR)

        Returns:
            Dict: 24-hour ticker data
        """
        endpoint = f"/v1/market/ticker24Hr/{contract_pair.lower()}"
        return self._client.get_request(endpoint, public=True)

    def get_agg_trades(self, contract_pair: str) -> Dict[str, Any]:
        """
        Get aggregated trade data for a specific trading pair.

        Args:
            contract_pair: Trading pair symbol (e.g., 'btc' for BTCINR)

        Returns:
            Dict: Aggregated trade data
        """
        endpoint = f"/v1/market/aggTrade/{contract_pair.lower()}"
        return self._client.get_request(endpoint, public=True)

    def get_depth(self, contract_pair: str) -> Dict[str, Any]:
        """
        Get order book depth data for a specific trading pair.

        Args:
            contract_pair: Trading pair symbol (e.g., 'btc' for BTCINR)

        Returns:
            Dict: Order book depth data
        """
        endpoint = f"/v1/market/depth/{contract_pair.lower()}"
        return self._client.get_request(endpoint, public=True)

    def get_klines(self, pair: str, interval: str, 
                  start_time: Optional[int] = None, 
                  end_time: Optional[int] = None, 
                  limit: Optional[int] = None) -> List[Dict[str, Any]]:
        """
        Get candlestick (kline) data for a specific trading pair and interval.

        Args:
            pair: Trading pair symbol (e.g., 'BTCINR')
            interval: Kline interval (e.g., '1m', '5m', '1h')
            start_time: Start time in milliseconds
            end_time: End time in milliseconds
            limit: Maximum number of klines to return

        Returns:
            List[Dict]: List of kline data
        """
        endpoint = "/v1/market/klines"
        
        params = {
            "pair": pair.upper(),
            "interval": interval.lower()
        }
        
        if start_time:
            params["startTime"] = start_time
        if end_time:
            params["endTime"] = end_time
        if limit:
            params["limit"] = limit
            
        return self._client.post_request(endpoint, params, public=True)
