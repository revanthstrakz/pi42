"""
Exchange API endpoints for Pi42.
"""

from typing import Dict, Any, Optional


class ExchangeAPI:
    """
    Exchange API client for managing Pi42 exchange settings.
    """

    def __init__(self, client):
        """
        Initialize the Exchange API client.

        Args:
            client: The main Pi42Client instance
        """
        self._client = client

    def get_exchange_info(self, market: Optional[str] = None) -> Dict[str, Any]:
        """
        Retrieve exchange information.

        Args:
            market: Market to return exchangeInfo for (e.g., 'INR' or 'USDT')

        Returns:
            Dict: Exchange information
        """
        endpoint = "/v1/exchange/exchangeInfo"
        
        params = {}
        if market:
            params["market"] = market
            
        return self._client.get_request(endpoint, params)

    def update_preference(self, leverage: int, margin_mode: str, 
                         contract_name: str) -> Dict[str, Any]:
        """
        Update the leverage and margin-mode for a specified contract.

        Args:
            leverage: The leverage level to set for the contract (must be an integer)
            margin_mode: Margin mode to set for the contract ('CROSS' or 'ISOLATED')
            contract_name: The trading pair or contract name (e.g., 'BTCINR')

        Returns:
            Dict: Preference update response
        """
        endpoint = "/v1/exchange/update/preference"
        
        params = {
            "leverage": leverage,
            "marginMode": margin_mode,
            "contractName": contract_name
        }
        
        return self._client.post_request(endpoint, params)

    def update_leverage(self, leverage: int, contract_name: str) -> Dict[str, Any]:
        """
        Update the leverage for a specified contract.

        Args:
            leverage: The leverage level to set for the contract (must be an integer)
            contract_name: The trading pair or contract name (e.g., 'BTCINR')

        Returns:
            Dict: Leverage update response
        """
        endpoint = "/v1/exchange/update/leverage"
        
        params = {
            "leverage": leverage,
            "contractName": contract_name
        }
        
        return self._client.post_request(endpoint, params)
