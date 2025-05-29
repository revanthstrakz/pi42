package pi42

import (
	"encoding/json"
	"fmt"
)

// ExchangeAPI provides access to exchange settings endpoints
type ExchangeAPI struct {
	client *Client
}

// NewExchangeAPI creates a new Exchange API instance
func NewExchangeAPI(client *Client) *ExchangeAPI {
	return &ExchangeAPI{client: client}
}

// ExchangeInfo retrieves exchange information with structured response
func (api *ExchangeAPI) ExchangeInfo(market string) (*ExchangeInfoResponse, error) {
	endpoint := "/v1/exchange/exchangeInfo"

	params := make(map[string]string)
	if market != "" {
		params["market"] = market
	}

	data, err := api.client.Get(endpoint, params, true)
	if err != nil {
		return nil, err
	}

	var result ExchangeInfoResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}

// UpdatePreference updates the leverage and margin-mode for a specified contract
func (api *ExchangeAPI) UpdatePreference(leverage int, marginMode, contractName string) (*PreferenceUpdateResponse, error) {
	endpoint := "/v1/exchange/update/preference"

	params := map[string]interface{}{
		"leverage":     leverage,
		"marginMode":   marginMode,
		"contractName": contractName,
	}

	data, err := api.client.Post(endpoint, params, false)
	if err != nil {
		return nil, err
	}

	var result PreferenceUpdateResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}

// UpdateLeverage updates the leverage for a specified contract
func (api *ExchangeAPI) UpdateLeverage(leverage int, contractName string) (*LeverageUpdateResponse, error) {
	endpoint := "/v1/exchange/update/leverage"

	params := map[string]interface{}{
		"leverage":     leverage,
		"contractName": contractName,
	}

	data, err := api.client.Post(endpoint, params, false)
	if err != nil {
		return nil, err
	}

	var result LeverageUpdateResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}
