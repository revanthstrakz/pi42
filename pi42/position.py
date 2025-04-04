"""
Position API endpoints for Pi42.
"""

from typing import Dict, Any, Optional, List, Union


class PositionAPI:
    """
    Position API client for managing Pi42 positions.
    """

    def __init__(self, client):
        """
        Initialize the Position API client.

        Args:
            client: The main Pi42Client instance
        """
        self._client = client

    def get_positions(self, position_status: str, 
                     start_timestamp: Optional[int] = None,
                     end_timestamp: Optional[int] = None,
                     sort_order: Optional[str] = None,
                     page_size: Optional[int] = None,
                     symbol: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Retrieve positions based on their status.

        Args:
            position_status: The status of positions to retrieve ('OPEN', 'CLOSED', 'LIQUIDATED')
            start_timestamp: Start timestamp for filtering positions
            end_timestamp: End timestamp for filtering positions
            sort_order: Sorting order ('asc' or 'desc')
            page_size: Number of results to return per page
            symbol: Trading pair symbol to filter positions

        Returns:
            List[Dict]: List of positions matching the criteria
        """
        endpoint = f"/v1/positions/{position_status.upper()}"
        
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

    def get_position(self, position_id: str) -> Dict[str, Any]:
        """
        Retrieve details for a specific position.

        Args:
            position_id: The unique identifier for the position to retrieve

        Returns:
            Dict: Position details
        """
        endpoint = f"/v1/positions"
        
        params = {
            "positionId": position_id
        }
        
        return self._client.get_request(endpoint, params)

    def close_all_positions(self) -> Dict[str, Any]:
        """
        Close all open positions.

        Returns:
            Dict: Position closure response
        """
        endpoint = "/v1/positions/close-all-positions"
        
        return self._client.delete_request(endpoint)
