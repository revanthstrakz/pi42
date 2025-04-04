"""
User Data API endpoints for Pi42.
"""

from typing import Dict, Any, Optional, List


class UserDataAPI:
    """
    User Data API client for accessing Pi42 user-specific data.
    """

    def __init__(self, client):
        """
        Initialize the User Data API client.

        Args:
            client: The main Pi42Client instance
        """
        self._client = client

    def get_trade_history(self, start_timestamp: Optional[int] = None,
                         end_timestamp: Optional[int] = None,
                         sort_order: Optional[str] = None,
                         page_size: Optional[int] = None,
                         symbol: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Retrieve the trade history for a user.

        Args:
            start_timestamp: Start timestamp for filtering trade history
            end_timestamp: End timestamp for filtering trade history
            sort_order: Sorting order ('asc' or 'desc')
            page_size: Number of results to return per page
            symbol: Trading pair symbol to filter trade history

        Returns:
            List[Dict]: List of trade history items
        """
        endpoint = "/v1/user-data/trade-history"
        
        params = {}
        
        if start_timestamp:
            params["startTimestamp"] = start_timestamp
        if end_timestamp:
            params["endTimestamp"] = end_timestamp
        if sort_order:
            params["sortOrder"] = sort_order
        if page_size:
            params["pageSize"] = page_size
        if symbol:
            params["symbol"] = symbol
            
        return self._client.get_request(endpoint, params)

    def get_transaction_history(self, start_timestamp: Optional[int] = None,
                               end_timestamp: Optional[int] = None,
                               sort_order: Optional[str] = None,
                               page_size: Optional[int] = None,
                               symbol: Optional[str] = None,
                               trade_id: Optional[int] = None,
                               position_id: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Retrieve the transaction history for a user.

        Args:
            start_timestamp: Start timestamp for filtering transaction history
            end_timestamp: End timestamp for filtering transaction history
            sort_order: Sorting order ('asc' or 'desc')
            page_size: Number of records to return per page
            symbol: Trading symbol to filter the transaction history
            trade_id: Specific trade ID to filter the transaction history
            position_id: Specific position ID to filter the transaction history

        Returns:
            List[Dict]: List of transaction history items
        """
        endpoint = "/v1/user-data/transaction-history"
        
        params = {}
        
        if start_timestamp:
            params["startTimestamp"] = start_timestamp
        if end_timestamp:
            params["endTimestamp"] = end_timestamp
        if sort_order:
            params["sortOrder"] = sort_order
        if page_size:
            params["pageSize"] = page_size
        if symbol:
            params["symbol"] = symbol
        if trade_id:
            params["tradeId"] = trade_id
        if position_id:
            params["positionId"] = position_id
            
        return self._client.get_request(endpoint, params)

    def create_listen_key(self) -> Dict[str, str]:
        """
        Create a new listen key for WebSocket connections.

        Returns:
            Dict: Response containing the listen key
        """
        endpoint = "/v1/retail/listen-key"
        
        return self._client.post_request(endpoint)

    def update_listen_key(self) -> str:
        """
        Update the listen key for WebSocket connections.

        Returns:
            str: Confirmation message
        """
        endpoint = "/v1/retail/listen-key"
        
        return self._client.put_request(endpoint)

    def delete_listen_key(self) -> str:
        """
        Delete the listen key for WebSocket connections.

        Returns:
            str: Confirmation message
        """
        endpoint = "/v1/retail/listen-key"
        
        return self._client.delete_request(endpoint)
