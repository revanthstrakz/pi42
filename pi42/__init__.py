"""
Pi42 API Client Package

This package provides a Python interface to the Pi42 API.
"""

from .client import Pi42Client
from .exceptions import Pi42APIError, Pi42RequestError

__all__ = ["Pi42Client", "Pi42APIError", "Pi42RequestError"]
