"""
Order API endpoints for Pi42.
"""

from typing import Dict, Any, Optional, List, Union


class OrderAPI:
    """
    Order API client for managing Pi42 orders.
    """

    def __init__(self, client):
        """
        Initialize the Order API client.

        Args:
            client: The main Pi42Client instance
        """
        self._client = client

    def place_order(self, symbol: str, side: str, order_type: str, quantity: float,
                   place_type: str = "ORDER_FORM", margin_asset: str = "INR",
                   price: Optional[float] = None, reduce_only: bool = False,
                   take_profit_price: Optional[float] = None,
                   stop_loss_price: Optional[float] = None,
                   stop_price: Optional[float] = None,
                   position_id: Optional[str] = None,
                   device_type: str = "WEB",
                   user_category: str = "EXTERNAL") -> Dict[str, Any]:
        """
        Place an order on Pi42's trading platform.

        Args:
            symbol: Trading pair symbol (e.g., 'BTCUSDT')
            side: Order side ('BUY' or 'SELL')
            order_type: Order type ('MARKET', 'LIMIT', 'STOP_MARKET', 'STOP_LIMIT')
            quantity: Amount of the asset to be ordered
            place_type: Type of order placement ('ORDER_FORM' or 'POSITION')
            margin_asset: Asset used for margin (e.g., 'INR')
            price: Price for LIMIT orders
            reduce_only: Whether the order should only reduce the position
            take_profit_price: Price at which take profit order should be executed
            stop_loss_price: Price at which stop loss order should be executed
            stop_price: Required for STOP_MARKET and STOP_LIMIT orders
            position_id: Position ID (required if placeType is 'POSITION')
            device_type: Device type ('WEB', 'MOBILE')
            user_category: User category ('EXTERNAL', 'INTERNAL')

        Returns:
            Dict: Order placement response
        """
        endpoint = "/v1/order/place-order"
        
        params = {
            "placeType": place_type,
            "quantity": quantity,
            "side": side,
            "symbol": symbol,
            "type": order_type,
            "reduceOnly": reduce_only,
            "marginAsset": margin_asset,
            "deviceType": device_type,
            "userCategory": user_category
        }
        
        # Add conditional parameters
        if price is not None:
            params["price"] = price
            
        if take_profit_price is not None:
            params["takeProfitPrice"] = take_profit_price
            
        if stop_loss_price is not None:
            params["stopLossPrice"] = stop_loss_price
            
        if stop_price is not None:
            params["stopPrice"] = stop_price
            
        if position_id is not None:
            params["positionId"] = position_id
            
        return self._client.post_request(endpoint, params)

    def add_margin(self, position_id: str, amount: Union[int, float]) -> Dict[str, Any]:
        """
        Add margin to a specific position.

        Args:
            position_id: Unique identifier for the position
            amount: Amount of margin to be added

        Returns:
            Dict: Response containing margin addition details
        """
        endpoint = "/v1/order/add-margin"
        
        params = {
            "positionId": position_id,
            "amount": amount
        }
        
        return self._client.post_request(endpoint, params)

    def reduce_margin(self, position_id: str, amount: Union[int, float]) -> Dict[str, Any]:
        """
        Reduce the margin on an existing trading position.

        Args:
            position_id: Unique identifier for the position
            amount: Amount of margin to reduce

        Returns:
            Dict: Response containing margin reduction details
        """
        endpoint = "/v1/order/reduce-margin"
        
        params = {
            "positionId": position_id,
            "amount": amount
        }
        
        return self._client.post_request(endpoint, params)

    def get_open_orders(self, page_size: Optional[int] = None, 
                       sort_order: Optional[str] = None,
                       start_timestamp: Optional[int] = None,
                       end_timestamp: Optional[int] = None,
                       symbol: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Retrieve open orders for the account.

        Args:
            page_size: Number of results to return per page
            sort_order: Sorting order ('asc' or 'desc')
            start_timestamp: Start timestamp for filtering open orders
            end_timestamp: End timestamp for filtering open orders
            symbol: Trading pair symbol to filter orders

        Returns:
            List[Dict]: List of open orders
        """
        endpoint = "/v1/order/open-orders"
        
        params = {}
        
        if page_size:
            params["pageSize"] = page_size
        if sort_order:
            params["sortOrder"] = sort_order
        if start_timestamp:
            params["startTimestamp"] = start_timestamp
        if end_timestamp:
            params["endTimestamp"] = end_timestamp
        if symbol:
            params["symbol"] = symbol
            
        return self._client.get_request(endpoint, params)

    def get_order_history(self, page_size: Optional[int] = None, 
                         sort_order: Optional[str] = None,
                         start_timestamp: Optional[int] = None,
                         end_timestamp: Optional[int] = None,
                         symbol: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Retrieve historical order data for the account.

        Args:
            page_size: Number of results to return per page
            sort_order: Sorting order ('asc' or 'desc')
            start_timestamp: Start timestamp for filtering order history
            end_timestamp: End timestamp for filtering order history
            symbol: Trading pair symbol to filter orders

        Returns:
            List[Dict]: List of historical orders
        """
        endpoint = "/v1/order/order-history"
        
        params = {}
        
        if page_size:
            params["pageSize"] = page_size
        if sort_order:
            params["sortOrder"] = sort_order
        if start_timestamp:
            params["startTimestamp"] = start_timestamp
        if end_timestamp:
            params["endTimestamp"] = end_timestamp
        if symbol:
            params["symbol"] = symbol
            
        return self._client.get_request(endpoint, params)

    def get_linked_orders(self, link_id: str) -> List[Dict[str, Any]]:
        """
        Retrieve orders that are linked by a specific link ID.

        Args:
            link_id: The unique identifier for the linked orders

        Returns:
            List[Dict]: List of linked orders
        """
        endpoint = f"/v1/order/linked-orders/{link_id}"
        
        return self._client.get_request(endpoint)

    def fetch_margin_history(self, symbol: Optional[str] = None,
                            page_size: Optional[int] = None, 
                            sort_order: Optional[str] = None,
                            start_timestamp: Optional[int] = None,
                            end_timestamp: Optional[int] = None) -> Dict[str, Any]:
        """
        Retrieve the margin history for an account.

        Args:
            symbol: Trading pair symbol to filter margin history
            page_size: Number of results to return per page
            sort_order: Sorting order ('asc' or 'desc')
            start_timestamp: Start timestamp for filtering margin history
            end_timestamp: End timestamp for filtering margin history

        Returns:
            Dict: Margin history response with data and pagination info
        """
        endpoint = "/v1/order/fetch-margin-history"
        
        params = {}
        
        if symbol:
            params["symbol"] = symbol
        if page_size:
            params["pageSize"] = page_size
        if sort_order:
            params["sortOrder"] = sort_order
        if start_timestamp:
            params["startTimestamp"] = start_timestamp
        if end_timestamp:
            params["endTimestamp"] = end_timestamp
            
        return self._client.get_request(endpoint, params)

    def delete_order(self, client_order_id: str) -> Dict[str, Any]:
        """
        Delete a specific order based on its client order ID.

        Args:
            client_order_id: The unique identifier for the order to be deleted

        Returns:
            Dict: Order deletion response
        """
        endpoint = "/v1/order/delete-order"
        
        params = {
            "clientOrderId": client_order_id
        }
        
        return self._client.delete_request(endpoint, params)

    def cancel_all_orders(self) -> Dict[str, Any]:
        """
        Cancel all open orders.

        Returns:
            Dict: Order cancellation response
        """
        endpoint = "/v1/order/cancel-all-orders"
        
        return self._client.delete_request(endpoint)
