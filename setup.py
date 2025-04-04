from setuptools import setup, find_packages

setup(
    name="pi42-api",
    version="0.1.0",
    description="Python client for Pi42 API",
    author="Pi42 API Team",
    packages=find_packages(),
    install_requires=[
        "requests>=2.25.0",
        "websocket-client>=1.0.0",
        "python-socketio>=5.0.0",
        "aiohttp>=3.7.4",
        "asyncio>=3.4.3",
    ],
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
    ],
    python_requires=">=3.7",
)
