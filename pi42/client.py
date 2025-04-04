"""
Main client class for the Pi42 API.
"""

import time
import json
import hmac
import hashlib
import requests
from typing import Dict, Any, Optional, Union, List

from .exceptions import Pi42APIError, Pi42RequestError
from .market import MarketAPI
from .order import OrderAPI
from .position import PositionAPI
from .wallet import WalletAPI
from .exchange import ExchangeAPI
from .user_data import UserDataAPI
from .websocket import WebSocketManager


class Pi42Client:
    """
    Pi42 API Client

    This client provides access to both authenticated and public endpoints of the Pi42 API.
    """

    def __init__(self, api_key: Optional[str] = None, api_secret: Optional[str] = None,
                 base_url: str = "https://fapi.pi42.com", public_url: str = "https://api.pi42.com"):
        """
        Initialize the Pi42 client.

        Args:
            api_key: Your Pi42 API key for authenticated endpoints
            api_secret: Your Pi42 API secret for request signing
            base_url: Base URL for authenticated endpoints
            public_url: Base URL for public endpoints
        """
        self.api_key = api_key
        self.api_secret = api_secret
        self.base_url = base_url
        self.public_url = public_url

        # Initialize API components
        self.market = MarketAPI(self)
        self.order = OrderAPI(self)
        self.position = PositionAPI(self)
        self.wallet = WalletAPI(self)
        self.exchange = ExchangeAPI(self)
        self.user_data = UserDataAPI(self)
        
        # Initialize WebSocket manager
        self.websocket = WebSocketManager(self)

    def _generate_signature(self, data_to_sign: str) -> str:
        """
        Generate HMAC SHA256 signature for request authentication.

        Args:
            data_to_sign: Data to be signed

        Returns:
            str: Hex-encoded signature
        """
        if not self.api_secret:
            raise ValueError("API secret is required for authenticated endpoints")
            
        return hmac.new(
            self.api_secret.encode('utf-8'),
            data_to_sign.encode('utf-8'),
            hashlib.sha256
        ).hexdigest()

    def _get_timestamp(self) -> str:
        """
        Get current timestamp in milliseconds.

        Returns:
            str: Current timestamp as a string
        """
        return str(int(time.time() * 1000))

    def _handle_response(self, response: requests.Response) -> Any:
        """
        Handle API response and raise appropriate exceptions if needed.

        Args:
            response: Response from the requests library

        Returns:
            dict: Parsed JSON response

        Raises:
            Pi42APIError: If the API returns an error
            Pi42RequestError: If there's an issue with the request
        """
        if not response.ok:
            try:
                error_data = response.json()
                error_message = error_data.get("message", "Unknown error")
                error_code = error_data.get("code")
                raise Pi42APIError(response.status_code, error_message, error_code)
            except (ValueError, KeyError):
                raise Pi42RequestError(f"HTTP Error: {response.status_code} - {response.text}")
        
        try:
            return response.json()
        except ValueError:
            return response.text

    def get_request(self, endpoint: str, params: Optional[Dict[str, Any]] = None, 
                    public: bool = False) -> Any:
        """
        Send a GET request to the Pi42 API.

        Args:
            endpoint: API endpoint to request
            params: Query parameters to include in the request
            public: Whether to use the public API URL

        Returns:
            Response data from the API

        Raises:
            Pi42APIError: If the API returns an error
            Pi42RequestError: If there's an issue with the request
        """
        base = self.public_url if public else self.base_url
        url = f"{base}{endpoint}"
        
        # Create a copy of the params dict or initialize an empty one
        actual_params = params.copy() if params else {}
        
        if not public:
            # Add timestamp for authenticated requests
            actual_params["timestamp"] = self._get_timestamp()
            
            # Create query string for signing
            query_string = "&".join([f"{k}={v}" for k, v in actual_params.items()])
            
            # Generate signature
            signature = self._generate_signature(query_string)
            
            # Set headers with API key and signature
            headers = {
                "api-key": self.api_key,
                "signature": signature,
                "accept": "*/*"
            }
        else:
            headers = {}
        
        try:
            response = requests.get(url, params=actual_params, headers=headers)
            return self._handle_response(response)
        except requests.RequestException as e:
            raise Pi42RequestError(str(e))

    def post_request(self, endpoint: str, params: Optional[Dict[str, Any]] = None, 
                     public: bool = False) -> Any:
        """
        Send a POST request to the Pi42 API.

        Args:
            endpoint: API endpoint to request
            params: Body parameters to include in the request
            public: Whether to use the public API URL

        Returns:
            Response data from the API

        Raises:
            Pi42APIError: If the API returns an error
            Pi42RequestError: If there's an issue with the request
        """
        base = self.public_url if public else self.base_url
        url = f"{base}{endpoint}"
        
        # Create a copy of the params dict or initialize an empty one
        actual_params = params.copy() if params else {}
        
        if not public:
            # Add timestamp for authenticated requests
            actual_params["timestamp"] = self._get_timestamp()
            
            # Convert the parameters to a JSON string for signing
            data_to_sign = json.dumps(actual_params, separators=(',', ':'))
            
            # Generate signature
            signature = self._generate_signature(data_to_sign)
            
            # Set headers with API key and signature
            headers = {
                "api-key": self.api_key,
                "Content-Type": "application/json",
                "signature": signature
            }
        else:
            headers = {
                "Content-Type": "application/json"
            }
        
        try:
            response = requests.post(url, json=actual_params, headers=headers)
            return self._handle_response(response)
        except requests.RequestException as e:
            raise Pi42RequestError(str(e))

    def put_request(self, endpoint: str, params: Optional[Dict[str, Any]] = None) -> Any:
        """
        Send a PUT request to the Pi42 API.

        Args:
            endpoint: API endpoint to request
            params: Body parameters to include in the request

        Returns:
            Response data from the API

        Raises:
            Pi42APIError: If the API returns an error
            Pi42RequestError: If there's an issue with the request
        """
        url = f"{self.base_url}{endpoint}"
        
        # Create a copy of the params dict or initialize an empty one
        actual_params = params.copy() if params else {}
        
        # Add timestamp for authenticated requests
        actual_params["timestamp"] = self._get_timestamp()
        
        # Convert the parameters to a JSON string for signing
        data_to_sign = json.dumps(actual_params, separators=(',', ':'))
        
        # Generate signature
        signature = self._generate_signature(data_to_sign)
        
        # Set headers with API key and signature
        headers = {
            "api-key": self.api_key,
            "Content-Type": "application/json",
            "signature": signature
        }
        
        try:
            response = requests.put(url, json=actual_params, headers=headers)
            return self._handle_response(response)
        except requests.RequestException as e:
            raise Pi42RequestError(str(e))

    def delete_request(self, endpoint: str, params: Optional[Dict[str, Any]] = None) -> Any:
        """
        Send a DELETE request to the Pi42 API.

        Args:
            endpoint: API endpoint to request
            params: Body parameters to include in the request

        Returns:
            Response data from the API

        Raises:
            Pi42APIError: If the API returns an error
            Pi42RequestError: If there's an issue with the request
        """
        url = f"{self.base_url}{endpoint}"
        
        # Create a copy of the params dict or initialize an empty one
        actual_params = params.copy() if params else {}
        
        # Add timestamp for authenticated requests
        actual_params["timestamp"] = self._get_timestamp()
        
        # Convert the parameters to a JSON string for signing
        data_to_sign = json.dumps(actual_params, separators=(',', ':'))
        
        # Generate signature
        signature = self._generate_signature(data_to_sign)
        
        # Set headers with API key and signature
        headers = {
            "api-key": self.api_key,
            "Content-Type": "application/json",
            "signature": signature
        }
        
        try:
            response = requests.delete(url, json=actual_params, headers=headers)
            return self._handle_response(response)
        except requests.RequestException as e:
            raise Pi42RequestError(str(e))
