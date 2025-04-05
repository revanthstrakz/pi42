package pi42

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// OrderAPI provides access to order management endpoints
type OrderAPI struct {
	client *Client
}

// NewOrderAPI creates a new Order API instance
func NewOrderAPI(client *Client) *OrderAPI {
	return &OrderAPI{client: client}
}

// PlaceOrderParams represents parameters for placing an order
type PlaceOrderParams struct {
	Symbol          string  `json:"symbol"`
	Side            string  `json:"side"`
	Type            string  `json:"type"`
	Quantity        float64 `json:"quantity"`
	PlaceType       string  `json:"placeType"`
	MarginAsset     string  `json:"marginAsset"`
	Price           float64 `json:"price,omitempty"`
	ReduceOnly      bool    `json:"reduceOnly"`
	TakeProfitPrice float64 `json:"takeProfitPrice,omitempty"`
	StopLossPrice   float64 `json:"stopLossPrice,omitempty"`
	StopPrice       float64 `json:"stopPrice,omitempty"`
	PositionID      string  `json:"positionId,omitempty"`
	DeviceType      string  `json:"deviceType"`
	UserCategory    string  `json:"userCategory"`
}

// PlaceOrder places an order on Pi42's trading platform
func (api *OrderAPI) PlaceOrder(params PlaceOrderParams) (map[string]interface{}, error) {
	endpoint := "/v1/order/place-order"

	// Convert struct to map for the request
	paramsMap := map[string]interface{}{
		"symbol":       params.Symbol,
		"side":         params.Side,
		"type":         params.Type,
		"quantity":     params.Quantity,
		"reduceOnly":   params.ReduceOnly,
		"marginAsset":  params.MarginAsset,
		"deviceType":   params.DeviceType,
		"userCategory": params.UserCategory,
	}

	if params.PlaceType != "" {
		paramsMap["placeType"] = params.PlaceType
	} else {
		paramsMap["placeType"] = "ORDER_FORM"
	}

	if params.Price != 0 {
		paramsMap["price"] = params.Price
	}

	if params.TakeProfitPrice != 0 {
		paramsMap["takeProfitPrice"] = params.TakeProfitPrice
	}

	if params.StopLossPrice != 0 {
		paramsMap["stopLossPrice"] = params.StopLossPrice
	}

	if params.StopPrice != 0 {
		paramsMap["stopPrice"] = params.StopPrice
	}

	if params.PositionID != "" {
		paramsMap["positionId"] = params.PositionID
	}

	data, err := api.client.Post(endpoint, paramsMap, false)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// AddMargin adds margin to a specific position
func (api *OrderAPI) AddMargin(positionID string, amount float64) (map[string]interface{}, error) {
	endpoint := "/v1/order/add-margin"

	params := map[string]interface{}{
		"positionId": positionID,
		"amount":     amount,
	}

	data, err := api.client.Post(endpoint, params, false)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// ReduceMargin reduces the margin on an existing trading position
func (api *OrderAPI) ReduceMargin(positionID string, amount float64) (map[string]interface{}, error) {
	endpoint := "/v1/order/reduce-margin"

	params := map[string]interface{}{
		"positionId": positionID,
		"amount":     amount,
	}

	data, err := api.client.Post(endpoint, params, false)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// OrderQueryParams represents parameters for querying orders
type OrderQueryParams struct {
	PageSize       int    `json:"pageSize,omitempty"`
	SortOrder      string `json:"sortOrder,omitempty"`
	StartTimestamp int64  `json:"startTimestamp,omitempty"`
	EndTimestamp   int64  `json:"endTimestamp,omitempty"`
	Symbol         string `json:"symbol,omitempty"`
}

// GetOpenOrders retrieves open orders for the account
func (api *OrderAPI) GetOpenOrders(params OrderQueryParams) ([]map[string]interface{}, error) {
	endpoint := "/v1/order/open-orders"

	queryParams := make(map[string]string)

	if params.PageSize > 0 {
		queryParams["pageSize"] = strconv.Itoa(params.PageSize)
	}
	if params.SortOrder != "" {
		queryParams["sortOrder"] = params.SortOrder
	}
	if params.StartTimestamp > 0 {
		queryParams["startTimestamp"] = strconv.FormatInt(params.StartTimestamp, 10)
	}
	if params.EndTimestamp > 0 {
		queryParams["endTimestamp"] = strconv.FormatInt(params.EndTimestamp, 10)
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

// GetOrderHistory retrieves historical order data for the account
func (api *OrderAPI) GetOrderHistory(params OrderQueryParams) ([]map[string]interface{}, error) {
	endpoint := "/v1/order/order-history"

	queryParams := make(map[string]string)

	if params.PageSize > 0 {
		queryParams["pageSize"] = strconv.Itoa(params.PageSize)
	}
	if params.SortOrder != "" {
		queryParams["sortOrder"] = params.SortOrder
	}
	if params.StartTimestamp > 0 {
		queryParams["startTimestamp"] = strconv.FormatInt(params.StartTimestamp, 10)
	}
	if params.EndTimestamp > 0 {
		queryParams["endTimestamp"] = strconv.FormatInt(params.EndTimestamp, 10)
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

// GetLinkedOrders retrieves orders that are linked by a specific link ID
func (api *OrderAPI) GetLinkedOrders(linkID string) ([]map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/v1/order/linked-orders/%s", linkID)

	data, err := api.client.Get(endpoint, nil, false)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// FetchMarginHistory retrieves the margin history for an account
func (api *OrderAPI) FetchMarginHistory(params OrderQueryParams) (map[string]interface{}, error) {
	endpoint := "/v1/order/fetch-margin-history"

	queryParams := make(map[string]string)

	if params.PageSize > 0 {
		queryParams["pageSize"] = strconv.Itoa(params.PageSize)
	}
	if params.SortOrder != "" {
		queryParams["sortOrder"] = params.SortOrder
	}
	if params.StartTimestamp > 0 {
		queryParams["startTimestamp"] = strconv.FormatInt(params.StartTimestamp, 10)
	}
	if params.EndTimestamp > 0 {
		queryParams["endTimestamp"] = strconv.FormatInt(params.EndTimestamp, 10)
	}
	if params.Symbol != "" {
		queryParams["symbol"] = params.Symbol
	}

	data, err := api.client.Get(endpoint, queryParams, false)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// DeleteOrder deletes a specific order based on its client order ID
func (api *OrderAPI) DeleteOrder(clientOrderID string) (map[string]interface{}, error) {
	endpoint := "/v1/order/delete-order"

	params := map[string]interface{}{
		"clientOrderId": clientOrderID,
	}

	data, err := api.client.Delete(endpoint, params)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// CancelAllOrders cancels all open orders
func (api *OrderAPI) CancelAllOrders() (map[string]interface{}, error) {
	endpoint := "/v1/order/cancel-all-orders"

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
