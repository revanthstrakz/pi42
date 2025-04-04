"""
Wallet API endpoints for Pi42.
"""

from typing import Dict, Any, Optional


class WalletAPI:
    """
    Wallet API client for accessing Pi42 wallet information.
    """

    def __init__(self, client):
        """
        Initialize the Wallet API client.

        Args:
            client: The main Pi42Client instance
        """
        self._client = client

    def get_futures_wallet_details(self, margin_asset: Optional[str] = "INR") -> Dict[str, Any]:
        """
        Get all details of Futures wallet.

        Args:
            margin_asset: Margin asset type (default: 'INR')

        Returns:
            Dict: Futures wallet details
        """
        endpoint = "/v1/wallet/futures-wallet/details"
        
        params = {}
        if margin_asset:
            params["marginAsset"] = margin_asset
            
        return self._client.get_request(endpoint, params)

    def get_funding_wallet_details(self, margin_asset: Optional[str] = "INR") -> Dict[str, Any]:
        """
        Get details of funding wallet.

        Args:
            margin_asset: Margin asset type (default: 'INR')

        Returns:
            Dict: Funding wallet details
        """
        endpoint = "/v1/wallet/funding-wallet/details"
        
        params = {}
        if margin_asset:
            params["marginAsset"] = margin_asset
            
        return self._client.get_request(endpoint, params)
