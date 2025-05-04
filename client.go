package pi42

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ContractInfo holds information about a trading contract/symbol
type ContractInfo struct {
	Symbol            string
	Name              string
	ContractName      string
	BaseAsset         string
	QuoteAsset        string
	PricePrecision    int
	QuantityPrecision int
	MinQuantity       float64
	MaxQuantity       float64
	MarketMinQuantity float64
	MarketMaxQuantity float64
	OrderTypes        []string
	MaxLeverage       float64
	MarginAssets      []string
	ContractType      string
	LiquidationFee    float64
	Tags              []string
}

// ExchangeInfoResponse represents the structure of the exchange info API response
type ExchangeInfoResponse struct {
	Markets   []string       `json:"markets"`
	Contracts []ContractData `json:"contracts"`
}

// ContractData represents the raw contract data from the API
type ContractData struct {
	Name                  string   `json:"name"`
	ContractName          string   `json:"contractName"`
	Slug                  string   `json:"slug"`
	Tags                  []string `json:"tags"`
	Filters               []Filter `json:"filters"`
	MakerFee              float64  `json:"makerFee"`
	TakerFee              float64  `json:"takerFee"`
	BaseAsset             string   `json:"baseAsset"`
	OrderTypes            []string `json:"orderTypes"`
	QuoteAsset            string   `json:"quoteAsset"`
	MaxLeverage           string   `json:"maxLeverage"`
	ContractType          string   `json:"contractType"`
	PricePrecision        string   `json:"pricePrecision"`
	QuantityPrecision     string   `json:"quantityPrecision"`
	MarginAssetsSupported []string `json:"marginAssetsSupported"`
}

// Filter represents trading filters applied to contracts
type Filter struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty,omitempty"`
	MaxQty     string `json:"maxQty,omitempty"`
	Limit      string `json:"limit,omitempty"`
	Notional   string `json:"notional,omitempty"`
}

// Client represents the API client for Pi42
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

	ExchangeInfo map[string]ContractInfo
}

// NewClient creates a new API client instance
func NewClient(apiKey, apiSecret string) *Client {
	client := &Client{
		APIKey:       apiKey,
		APISecret:    apiSecret,
		BaseURL:      "https://fapi.pi42.com",
		PublicURL:    "https://api.pi42.com",
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		ExchangeInfo: make(map[string]ContractInfo),
	}

	// Initialize API components
	client.Market = NewMarketAPI(client)
	client.Order = NewOrderAPI(client)
	client.Position = NewPositionAPI(client)
	client.Wallet = NewWalletAPI(client)
	client.Exchange = NewExchangeAPI(client)
	client.UserData = NewUserDataAPI(client)
	err := client.fetchExchangeInfo()
	if err != nil {
		log.Printf("Error fetching exchange info: %v", err)
	} else {
		log.Println("Exchange info loaded successfully")
	}
	return client
}

// fetchExchangeInfo loads contract specifications from the exchange
func (c *Client) fetchExchangeInfo() error {
	endpoint := "/v1/exchange/exchangeInfo"

	data, err := c.Get(endpoint, nil, true)
	if err != nil {
		return err
	}

	var response ExchangeInfoResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("error parsing exchange info response: %v", err)
	}

	// Process each contract and extract the needed information
	for _, contract := range response.Contracts {
		// Parse precision values
		pricePrecision, _ := strconv.Atoi(contract.PricePrecision)
		quantityPrecision, _ := strconv.Atoi(contract.QuantityPrecision)
		maxLeverage, _ := strconv.ParseFloat(contract.MaxLeverage, 64)

		// Initialize with defaults
		contractInfo := ContractInfo{
			Symbol:            contract.Name,
			Name:              contract.Name,
			ContractName:      contract.ContractName,
			BaseAsset:         contract.BaseAsset,
			QuoteAsset:        contract.QuoteAsset,
			PricePrecision:    pricePrecision,
			QuantityPrecision: quantityPrecision,
			OrderTypes:        contract.OrderTypes,
			MaxLeverage:       maxLeverage,
			MarginAssets:      contract.MarginAssetsSupported,
			ContractType:      contract.ContractType,
			Tags:              contract.Tags,
		}

		// Extract filter information
		for _, filter := range contract.Filters {
			switch filter.FilterType {
			case "LIMIT_QTY_SIZE":
				contractInfo.MinQuantity, _ = strconv.ParseFloat(filter.MinQty, 64)
				contractInfo.MaxQuantity, _ = strconv.ParseFloat(filter.MaxQty, 64)
			case "MARKET_QTY_SIZE":
				contractInfo.MarketMinQuantity, _ = strconv.ParseFloat(filter.MinQty, 64)
				contractInfo.MarketMaxQuantity, _ = strconv.ParseFloat(filter.MaxQty, 64)
			}
		}
		c.ExchangeInfo[contract.Name] = contractInfo
	}

	return nil
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
