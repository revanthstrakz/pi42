package pi42

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Client represents the Pi42 API client
type Client struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	PublicURL  string
	HTTPClient *http.Client

	Market   *MarketAPI
	Order    *OrderAPI
	Position *PositionAPI
	Wallet   *WalletAPI
	Exchange *ExchangeAPI
	UserData *UserDataAPI
	Socketio *SocketioManager
}

// NewClient creates a new Pi42 API client
func NewClient(apiKey, apiSecret string) *Client {
	client := &Client{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		BaseURL:    "https://fapi.pi42.com",
		PublicURL:  "https://api.pi42.com",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}

	// Initialize API components
	client.Market = NewMarketAPI(client)
	client.Order = NewOrderAPI(client)
	client.Position = NewPositionAPI(client)
	client.Wallet = NewWalletAPI(client)
	client.Exchange = NewExchangeAPI(client)
	client.UserData = NewUserDataAPI(client)
	client.Socketio = NewSocketioManager(client)

	return client
}

// generateSignature creates an HMAC SHA256 signature for request authentication
func (c *Client) generateSignature(data string) (string, error) {
	if c.APISecret == "" {
		return "", fmt.Errorf("API secret is required for authenticated endpoints")
	}

	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}

// getTimestamp returns the current timestamp in milliseconds
func (c *Client) getTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

// Get sends a GET request to the Pi42 API
func (c *Client) Get(endpoint string, params map[string]string, public bool) ([]byte, error) {
	baseURL := c.PublicURL
	if !public {
		baseURL = c.BaseURL
	}

	// Build the URL
	requestURL := fmt.Sprintf("%s%s", baseURL, endpoint)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add query parameters
	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}

	// For authenticated requests, add timestamp and signature
	if !public {
		timestamp := c.getTimestamp()
		q.Add("timestamp", timestamp)

		// Create the query string for signing
		queryString := q.Encode()
		signature, err := c.generateSignature(queryString)
		if err != nil {
			return nil, err
		}

		// Add headers for authentication
		req.Header.Add("api-key", c.APIKey)
		req.Header.Add("signature", signature)
		req.Header.Add("accept", "*/*")
	}

	// Set the query parameters
	req.URL.RawQuery = q.Encode()

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for error responses - add special handling for 201 Created status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var apiError APIError
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, apiError
		}
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Post sends a POST request to the Pi42 API
func (c *Client) Post(endpoint string, params map[string]interface{}, public bool) ([]byte, error) {
	baseURL := c.PublicURL
	if !public {
		baseURL = c.BaseURL
	}

	// Add timestamp for authenticated requests
	if !public {
		params["timestamp"] = c.getTimestamp()
	}

	// Convert params to JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Build the URL
	requestURL := fmt.Sprintf("%s%s", baseURL, endpoint)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set content type header
	req.Header.Add("Content-Type", "application/json")

	// For authenticated requests, generate and add signature
	if !public {
		signature, err := c.generateSignature(string(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Add("api-key", c.APIKey)
		req.Header.Add("signature", signature)
	}

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for error responses - add special handling for 201 Created status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var apiError APIError
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, apiError
		}
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Put sends a PUT request to the Pi42 API
func (c *Client) Put(endpoint string, params map[string]interface{}) ([]byte, error) {
	// Add timestamp for authenticated requests
	params["timestamp"] = c.getTimestamp()

	// Convert params to JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Build the URL
	requestURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set content type header
	req.Header.Add("Content-Type", "application/json")

	// Generate and add signature
	signature, err := c.generateSignature(string(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-key", c.APIKey)
	req.Header.Add("signature", signature)

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for error responses
	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, apiError
		}
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Delete sends a DELETE request to the Pi42 API
func (c *Client) Delete(endpoint string, params map[string]interface{}) ([]byte, error) {
	// Add timestamp for authenticated requests
	params["timestamp"] = c.getTimestamp()

	// Convert params to JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Build the URL
	requestURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("DELETE", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set content type header
	req.Header.Add("Content-Type", "application/json")

	// Generate and add signature
	signature, err := c.generateSignature(string(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-key", c.APIKey)
	req.Header.Add("signature", signature)

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for error responses
	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, apiError
		}
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}
