package pi42

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// UserDataAPI provides access to user-specific data endpoints
type UserDataAPI struct {
	client *Client
}

// NewUserDataAPI creates a new User Data API instance
func NewUserDataAPI(client *Client) *UserDataAPI {
	return &UserDataAPI{client: client}
}

// DataQueryParams represents common parameters for data queries
type DataQueryParams struct {
	StartTimestamp int64  `json:"startTimestamp,omitempty"`
	EndTimestamp   int64  `json:"endTimestamp,omitempty"`
	SortOrder      string `json:"sortOrder,omitempty"`
	PageSize       int    `json:"pageSize,omitempty"`
	Symbol         string `json:"symbol,omitempty"`
}

// GetTradeHistory retrieves the trade history for a user
func (api *UserDataAPI) GetTradeHistory(params DataQueryParams) ([]map[string]interface{}, error) {
	endpoint := "/v1/user-data/trade-history"

	queryParams := make(map[string]string)

	if params.StartTimestamp > 0 {
		queryParams["startTimestamp"] = strconv.FormatInt(params.StartTimestamp, 10)
	}
	if params.EndTimestamp > 0 {
		queryParams["endTimestamp"] = strconv.FormatInt(params.EndTimestamp, 10)
	}
	if params.SortOrder != "" {
		queryParams["sortOrder"] = params.SortOrder
	}
	if params.PageSize > 0 {
		queryParams["pageSize"] = strconv.Itoa(params.PageSize)
	}
	if params.Symbol != "" {
		queryParams["symbol"] = params.Symbol
	}

	data, err := api.client.Get(endpoint, queryParams, false)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// TransactionHistoryParams extends DataQueryParams with additional fields
type TransactionHistoryParams struct {
	DataQueryParams
	TradeID    int    `json:"tradeId,omitempty"`
	PositionID string `json:"positionId,omitempty"`
}

// GetTransactionHistory retrieves the transaction history for a user
func (api *UserDataAPI) GetTransactionHistory(params TransactionHistoryParams) ([]map[string]interface{}, error) {
	endpoint := "/v1/user-data/transaction-history"

	queryParams := make(map[string]string)

	if params.StartTimestamp > 0 {
		queryParams["startTimestamp"] = strconv.FormatInt(params.StartTimestamp, 10)
	}
	if params.EndTimestamp > 0 {
		queryParams["endTimestamp"] = strconv.FormatInt(params.EndTimestamp, 10)
	}
	if params.SortOrder != "" {
		queryParams["sortOrder"] = params.SortOrder
	}
	if params.PageSize > 0 {
		queryParams["pageSize"] = strconv.Itoa(params.PageSize)
	}
	if params.Symbol != "" {
		queryParams["symbol"] = params.Symbol
	}
	if params.TradeID > 0 {
		queryParams["tradeId"] = strconv.Itoa(params.TradeID)
	}
	if params.PositionID != "" {
		queryParams["positionId"] = params.PositionID
	}

	data, err := api.client.Get(endpoint, queryParams, false)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// CreateListenKey creates a new listen key for Socketio connections
func (api *UserDataAPI) CreateListenKey() (map[string]string, error) {
	endpoint := "/v1/retail/listen-key"

	data, err := api.client.Post(endpoint, map[string]interface{}{}, false)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// UpdateListenKey updates the listen key for Socketio connections
func (api *UserDataAPI) UpdateListenKey() (string, error) {
	endpoint := "/v1/retail/listen-key"

	data, err := api.client.Put(endpoint, map[string]interface{}{})
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DeleteListenKey deletes the listen key for Socketio connections
func (api *UserDataAPI) DeleteListenKey() (string, error) {
	endpoint := "/v1/retail/listen-key"

	data, err := api.client.Delete(endpoint, map[string]interface{}{})
	if err != nil {
		return "", err
	}

	return string(data), nil
}
