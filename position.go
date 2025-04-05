package pi42

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// PositionAPI provides access to position management endpoints
type PositionAPI struct {
	client *Client
}

// NewPositionAPI creates a new Position API instance
func NewPositionAPI(client *Client) *PositionAPI {
	return &PositionAPI{client: client}
}

// PositionQueryParams represents parameters for querying positions
type PositionQueryParams struct {
	StartTimestamp int64  `json:"startTimestamp,omitempty"`
	EndTimestamp   int64  `json:"endTimestamp,omitempty"`
	SortOrder      string `json:"sortOrder,omitempty"`
	PageSize       int    `json:"pageSize,omitempty"`
	Symbol         string `json:"symbol,omitempty"`
}

// GetPositions retrieves positions based on their status
func (api *PositionAPI) GetPositions(positionStatus string, params PositionQueryParams) ([]map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/v1/positions/%s", strings.ToUpper(positionStatus))

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

// GetPosition retrieves details for a specific position
func (api *PositionAPI) GetPosition(positionID string) (map[string]interface{}, error) {
	endpoint := "/v1/positions"

	params := map[string]string{
		"positionId": positionID,
	}

	data, err := api.client.Get(endpoint, params, false)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// CloseAllPositions closes all open positions
func (api *PositionAPI) CloseAllPositions() (map[string]interface{}, error) {
	endpoint := "/v1/positions/close-all-positions"

	data, err := api.client.Delete(endpoint, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}
