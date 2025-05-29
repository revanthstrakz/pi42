package pi42

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MarketAPI provides access to market data endpoints
type MarketAPI struct {
	client *Client
}

// NewMarketAPI creates a new Market API instance
func NewMarketAPI(client *Client) *MarketAPI {
	return &MarketAPI{client: client}
}

// GetTicker24hr gets 24-hour ticker data for a specific trading pair
func (api *MarketAPI) GetTicker24hr(contractPair string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/v1/market/ticker24Hr/%s", strings.ToLower(contractPair))
	params := make(map[string]string)

	data, err := api.client.Get(endpoint, params, true)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// GetAggTrades gets aggregated trade data for a specific trading pair
func (api *MarketAPI) GetAggTrades(contractPair string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/v1/market/aggTrade/%s", strings.ToLower(contractPair))

	data, err := api.client.Get(endpoint, nil, true)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// GetDepth gets order book depth data for a specific trading pair
// Returns structured DepthResponse containing order book bids and asks
func (api *MarketAPI) GetDepth(contractPair string) (*DepthResponse, error) {
	endpoint := fmt.Sprintf("/v1/market/depth/%s", strings.ToLower(contractPair))

	data, err := api.client.Get(endpoint, nil, true)
	if err != nil {
		return nil, err
	}

	var result DepthResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing depth response: %v", err)
	}

	return &result, nil
}

// KlinesParams represents parameters for the Klines method
type KlinesParams struct {
	Pair      string `json:"pair"`                // Trading pair (e.g., "BTCINR")
	Interval  string `json:"interval"`            // Kline interval (e.g., "1m", "1h", "1d")
	StartTime int64  `json:"startTime,omitempty"` // Optional start time in milliseconds
	EndTime   int64  `json:"endTime,omitempty"`   // Optional end time in milliseconds
	Limit     int    `json:"limit,omitempty"`     // Optional limit on number of results
}

// GetKlines gets candlestick (kline) data for a specific trading pair and interval
// Returns an array of structured KlineData objects
func (api *MarketAPI) GetKlines(params KlinesParams) ([]KlineData, error) {
	endpoint := "/v1/market/klines"

	// Convert struct to map for the request
	paramsMap := map[string]interface{}{
		"pair":     strings.ToUpper(params.Pair),
		"interval": strings.ToLower(params.Interval),
	}

	if params.StartTime > 0 {
		paramsMap["startTime"] = params.StartTime
	}
	if params.EndTime > 0 {
		paramsMap["endTime"] = params.EndTime
	}
	if params.Limit > 0 {
		paramsMap["limit"] = params.Limit
	}

	data, err := api.client.Post(endpoint, paramsMap, true)
	if err != nil {
		return nil, err
	}

	var result []KlineData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing klines response: %v", err)
	}

	return result, nil
}

// For backward compatibility
func (api *MarketAPI) Ticker24Hr(contractPair string) (map[string]interface{}, error) {
	return api.GetTicker24hr(contractPair)
}