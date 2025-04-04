"""
Exceptions for the Pi42 API Client.
"""

class Pi42APIError(Exception):
    """Exception raised for API errors from the Pi42 server."""
    def __init__(self, status_code, message, error_code=None):
        self.status_code = status_code
        self.message = message
        self.error_code = error_code
        super().__init__(f"API Error (Code: {error_code}, Status: {status_code}): {message}")


class Pi42RequestError(Exception):
    """Exception raised for network or request errors."""
    def __init__(self, message):
        self.message = message
        super().__init__(f"Request Error: {message}")
