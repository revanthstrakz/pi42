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

	// Use proper params map instead of nil
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
func (api *MarketAPI) GetDepth(contractPair string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/v1/market/depth/%s", strings.ToLower(contractPair))

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

// KlinesParams represents parameters for the Klines method
type KlinesParams struct {
	Pair      string `json:"pair"`
	Interval  string `json:"interval"`
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// GetKlines gets candlestick (kline) data for a specific trading pair and interval
func (api *MarketAPI) GetKlines(params KlinesParams) ([]map[string]interface{}, error) {
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

	// Fix the parsing for the klines response
	// The API seems to be returning valid data but it's treated as an error
	// Let's try to parse it as a successful response even if status code isn't 200
	var result []map[string]interface{}

	// Check if the data starts with '[' which indicates a JSON array
	if len(data) > 0 && data[0] == '[' {
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error parsing klines response: %v", err)
		}
		return result, nil
	}

	// Otherwise try to parse as a normal JSON object
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// For backward compatibility
func (api *MarketAPI) Ticker24Hr(contractPair string) (map[string]interface{}, error) {
	return api.GetTicker24hr(contractPair)
}

// For backward compatibility
func (api *MarketAPI) Klines(params KlinesParams) ([]map[string]interface{}, error) {
	return api.GetKlines(params)
}
