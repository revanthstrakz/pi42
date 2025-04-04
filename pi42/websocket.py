"""
WebSocket implementation for Pi42 API.
"""

import json
import socketio
import asyncio
from typing import Dict, Any, Optional, List, Callable, Union


class WebSocketManager:
    """
    WebSocket manager for Pi42 API.
    """

    def __init__(self, client):
        """
        Initialize the WebSocket manager.

        Args:
            client: The main Pi42Client instance
        """
        self._client = client
        self._public_sio = None
        self._auth_sio = None
        self._public_url = "https://fawss.pi42.com/"
        self._auth_url = "https://fawss-uds.pi42.com/auth-stream"
        self._listen_key = None
        self._callbacks = {}

    async def connect_public(self, topics: List[str] = None):
        """
        Connect to the public WebSocket server and subscribe to topics.

        Args:
            topics: List of topics to subscribe to
        """
        self._public_sio = socketio.AsyncClient()
        
        @self._public_sio.event
        async def connect():
            print('Connected to Public WebSocket server')
            if topics:
                await self.subscribe_public(topics)
        
        @self._public_sio.event
        async def disconnect():
            print('Disconnected from Public WebSocket server')
        
        # Register event handlers for different data types
        for event_type in [
            'depthUpdate', 'kline', 'markPriceUpdate', 'aggTrade', 
            '24hrTicker', 'marketInfo', 'markPriceArr', 'tickerArr'
        ]:
            self._register_public_event_handler(event_type)
        
        await self._public_sio.connect(self._public_url, transports=["websocket"])

    def _register_public_event_handler(self, event_type: str):
        """
        Register an event handler for a public WebSocket event.

        Args:
            event_type: The type of event to handle
        """
        @self._public_sio.on(event_type)
        async def handler(data):
            callback = self._callbacks.get(event_type)
            if callback:
                callback(data)
            else:
                print(f"{event_type}:", data)

    async def subscribe_public(self, topics: List[str]):
        """
        Subscribe to public WebSocket topics.

        Args:
            topics: List of topics to subscribe to
        """
        if not self._public_sio:
            raise ValueError("Not connected to public WebSocket server")
            
        await self._public_sio.emit('subscribe', {
            'params': topics
        })
        print(f"Subscribed to {topics}")

    async def connect_authenticated(self, listen_key: Optional[str] = None):
        """
        Connect to the authenticated WebSocket server.

        Args:
            listen_key: Listen key for authentication (if not provided, will try to create one)
        """
        if not listen_key:
            if not self._client.api_key or not self._client.api_secret:
                raise ValueError("API key and secret are required for authenticated WebSocket")
                
            # Create a listen key if not provided
            response = self._client.user_data.create_listen_key()
            listen_key = response.get("listenKey")
            
        self._listen_key = listen_key
        self._auth_sio = socketio.AsyncClient()
        namespace = f"/auth-stream/{listen_key}"
        
        @self._auth_sio.event
        async def connect():
            print('Connected to Authenticated WebSocket server')
        
        @self._auth_sio.event
        async def disconnect():
            print('Disconnected from Authenticated WebSocket server')
        
        # Register event handlers for authenticated events
        for event_type in [
            'newPosition', 'orderFilled', 'orderPartiallyFilled', 'orderCancelled', 
            'orderFailed', 'newOrder', 'updateOrder', 'updatePosition', 
            'closePosition', 'balanceUpdate', 'newTrade', 'sessionExpired'
        ]:
            self._register_auth_event_handler(event_type, namespace)
        
        await self._auth_sio.connect(self._auth_url, transports=["websocket"])

    def _register_auth_event_handler(self, event_type: str, namespace: str):
        """
        Register an event handler for an authenticated WebSocket event.

        Args:
            event_type: The type of event to handle
            namespace: The namespace for the authenticated events
        """
        @self._auth_sio.on(event_type, namespace=namespace)
        async def handler(data):
            callback = self._callbacks.get(event_type)
            if callback:
                callback(data)
            else:
                print(f"{event_type}:", data)

    def on(self, event_type: str, callback: Callable[[Dict[str, Any]], None]):
        """
        Register a callback for a specific event type.

        Args:
            event_type: The type of event to handle
            callback: Function to call when the event is received
        """
        self._callbacks[event_type] = callback

    async def close(self):
        """
        Close all WebSocket connections.
        """
        tasks = []
        
        if self._public_sio and self._public_sio.connected:
            tasks.append(self._public_sio.disconnect())
            
        if self._auth_sio and self._auth_sio.connected:
            tasks.append(self._auth_sio.disconnect())
            
        if tasks:
            await asyncio.gather(*tasks)
            
        self._public_sio = None
        self._auth_sio = None


# Create a simpler synchronous wrapper for the WebSocket manager
class WebSocketClient:
    """
    Synchronous wrapper for WebSocketManager.
    """

    def __init__(self, client):
        """
        Initialize the WebSocket client.

        Args:
            client: The main Pi42Client instance
        """
        self._manager = WebSocketManager(client)
        self._loop = None
        self._task = None

    def connect_public(self, topics: List[str] = None):
        """
        Connect to the public WebSocket server and subscribe to topics.

        Args:
            topics: List of topics to subscribe to
        """
        self._loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self._loop)
        
        async def _connect():
            await self._manager.connect_public(topics)
            
        self._task = self._loop.create_task(_connect())
        
        # Run the event loop in a separate thread
        import threading
        threading.Thread(target=self._loop.run_forever, daemon=True).start()

    def connect_authenticated(self, listen_key: Optional[str] = None):
        """
        Connect to the authenticated WebSocket server.

        Args:
            listen_key: Listen key for authentication
        """
        if not self._loop:
            self._loop = asyncio.new_event_loop()
            asyncio.set_event_loop(self._loop)
            
            # Run the event loop in a separate thread
            import threading
            threading.Thread(target=self._loop.run_forever, daemon=True).start()
        
        async def _connect():
            await self._manager.connect_authenticated(listen_key)
            
        self._loop.call_soon_threadsafe(lambda: self._loop.create_task(_connect()))

    def on(self, event_type: str, callback: Callable[[Dict[str, Any]], None]):
        """
        Register a callback for a specific event type.

        Args:
            event_type: The type of event to handle
            callback: Function to call when the event is received
        """
        self._manager.on(event_type, callback)

    def close(self):
        """
        Close all WebSocket connections.
        """
        if self._loop:
            async def _close():
                await self._manager.close()
                
            future = asyncio.run_coroutine_threadsafe(_close(), self._loop)
            future.result()  # Wait for completion
            
            self._loop.call_soon_threadsafe(self._loop.stop)
            self._loop = None
            self._task = None
