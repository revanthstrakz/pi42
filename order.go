package pi42

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	Symbol          string    `json:"symbol"`
	Side            OrderSide `json:"side"`
	Type            OrderType `json:"type"`
	Quantity        float64   `json:"quantity"`
	PlaceType       string    `json:"placeType"`
	MarginAsset     string    `json:"marginAsset"`
	Price           float64   `json:"price,omitempty"`
	ReduceOnly      bool      `json:"reduceOnly"`
	TakeProfitPrice float64   `json:"takeProfitPrice,omitempty"`
	StopLossPrice   float64   `json:"stopLossPrice,omitempty"`
	StopPrice       float64   `json:"stopPrice,omitempty"`
	PositionID      string    `json:"positionId,omitempty"`
	DeviceType      string    `json:"deviceType"`
	UserCategory    string    `json:"userCategory"`
	Leverage        int       `json:"leverage,omitempty"`
}

// OrderResponse represents the structured response when placing an order
type OrderResponse struct {
	ClientOrderID       string  `json:"clientOrderId"`
	Time                string  `json:"time"`
	Symbol              string  `json:"symbol"`
	ContractType        string  `json:"contractType"`
	Type                string  `json:"type"`
	Side                string  `json:"side"`
	Price               float64 `json:"price"`
	OrderAmount         float64 `json:"orderAmount"`
	FilledAmount        float64 `json:"filledAmount"`
	AvailableBalance    float64 `json:"availableBalance"`
	LinkID              string  `json:"linkId"`
	LinkType            string  `json:"linkType"`
	SubType             string  `json:"subType"`
	PlaceType           string  `json:"placeType"`
	LockedMargin        float64 `json:"lockedMargin"`
	BaseAsset           string  `json:"baseAsset"`
	QuoteAsset          string  `json:"quoteAsset"`
	MarginAsset         string  `json:"marginAsset"`
	LockedMarginInAsset float64 `json:"lockedMarginInMarginAsset"`
	Leverage            int     `json:"leverage"`
	ID                  float64 `json:"id"`
	StopPrice           float64 `json:"stopPrice"`
}

// PlaceOrder places an order on Pi42's trading platform
func (api *OrderAPI) PlaceOrder(params PlaceOrderParams) (OrderResponse, error) {
	endpoint := "/v1/order/place-order"

	// Convert struct to map for the request
	paramsMap := map[string]interface{}{
		"symbol":      params.Symbol,
		"side":        params.Side,
		"type":        params.Type,
		"quantity":    params.Quantity,
		"reduceOnly":  params.ReduceOnly,
		"marginAsset": params.MarginAsset,
	}

	// Add other parameters if they're set
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
		paramsMap["placeType"] = "POSITION"
	} else {
		paramsMap["placeType"] = "ORDER_FORM"
	}

	if params.Leverage > 0 {
		paramsMap["leverage"] = params.Leverage
	}

	data, err := api.client.Post(endpoint, paramsMap, false)
	if err != nil {
		return OrderResponse{}, err
	}

	var result OrderResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return OrderResponse{}, fmt.Errorf("error parsing response: %v", err)
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

// GetOpenOrders retrieves open orders for the account with structured response
func (api *OrderAPI) GetOpenOrders(params OrderQueryParams) ([]OpenOrder, error) {
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

	var result []OpenOrder
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// GetOrderHistory retrieves historical order data with structured response
func (api *OrderAPI) GetOrderHistory(params OrderQueryParams) ([]OrderHistoryItem, error) {
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

	var result []OrderHistoryItem
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// GetLinkedOrders retrieves orders that are linked by a specific link ID
func (api *OrderAPI) GetLinkedOrders(linkID string) ([]LinkedOrder, error) {
	endpoint := fmt.Sprintf("/v1/order/linked-orders/%s", linkID)

	data, err := api.client.Get(endpoint, nil, false)
	if err != nil {
		return nil, err
	}

	var result []LinkedOrder
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
func (api *OrderAPI) DeleteOrder(clientOrderID string) (*OrderCancelResponse, error) {
	endpoint := "/v1/order/delete-order"

	params := map[string]interface{}{
		"clientOrderId": clientOrderID,
	}

	data, err := api.client.Delete(endpoint, params)
	if err != nil {
		return nil, err
	}

	var result OrderCancelResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}

// CancelAllOrders cancels all open orders with structured response
func (api *OrderAPI) CancelAllOrders() (*BatchCancelResponse, error) {
	endpoint := "/v1/order/cancel-all-orders"

	data, err := api.client.Delete(endpoint, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var result BatchCancelResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}

// BulletParams represents simplified parameters for quick order placement
type BulletParams struct {
	Symbol     string    // Trading pair symbol
	Side       OrderSide // BUY or SELL
	OrderType  OrderType // MARKET, LIMIT, STOP_MARKET, or STOP_LIMIT
	Price      float64   // Required for LIMIT and STOP_LIMIT orders
	StopPrice  float64   // Required for STOP_MARKET and STOP_LIMIT orders
	Count      float64   // Multiplier for minimum quantity
	ReduceOnly bool      // Whether this is a reduce-only order
	Leverage   int       // Leverage to use for the order (optional)
	PositionID string    // Position ID for the order (optional)
}

// Bullet creates an order using exchange specifications for precision and minimum quantity
// and returns a structured order response
func (api *OrderAPI) Bullet(params BulletParams) (*OrderResponse, error) {
	// Get contract info for the symbol
	contractInfo, ok := api.client.ExchangeInfo[params.Symbol]
	if !ok {
		return nil, fmt.Errorf("symbol %s not found in exchange info", params.Symbol)
	}

	// Validate order type
	validOrderTypes := []OrderType{"MARKET", "LIMIT", "STOP_MARKET", "STOP_LIMIT"}
	isValidType := false
	for _, orderType := range validOrderTypes {
		if params.OrderType == orderType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return nil, fmt.Errorf("invalid order type: %s. Must be MARKET, LIMIT, STOP_MARKET, or STOP_LIMIT", params.OrderType)
	}

	// Validate required parameters for specific order types
	if (params.OrderType == "LIMIT" || params.OrderType == "STOP_LIMIT") && params.Price <= 0 {
		return nil, fmt.Errorf("price must be specified and greater than 0 for %s orders", params.OrderType)
	}

	if (params.OrderType == "STOP_MARKET" || params.OrderType == "STOP_LIMIT") && params.StopPrice <= 0 {
		return nil, fmt.Errorf("stopPrice must be specified and greater than 0 for %s orders", params.OrderType)
	}

	// Determine the minimum quantity based on order type
	var minQuantity float64
	var maxQuantity float64

	if params.OrderType == "MARKET" || params.OrderType == "STOP_MARKET" {
		minQuantity = contractInfo.MarketMinQuantity
		maxQuantity = contractInfo.MarketMaxQuantity
	} else {
		minQuantity = contractInfo.MinQuantity
		maxQuantity = contractInfo.MaxQuantity
	}

	// If min quantity is not set (could happen if filter parsing failed), use a safe default
	if minQuantity <= 0 {
		log.Default().Printf("Warning: Minimum quantity for %s not set, using default value\n", params.Symbol)
		minQuantity = 0.001 // Default fallback
	}

	// Calculate quantity based on minimum quantity and count
	quantity := minQuantity * params.Count

	// Check if quantity exceeds the maximum
	if maxQuantity > 0 && quantity > maxQuantity {
		return nil, fmt.Errorf("calculated quantity %.8f exceeds maximum allowed %.8f for %s",
			quantity, maxQuantity, params.Symbol)
	}

	// Round to the correct precision
	quantity = roundToDecimal(quantity, contractInfo.QuantityPrecision)

	// Check if the order type is supported for this symbol
	// For stop orders, we check if the base type (MARKET/LIMIT) is supported
	baseOrderType := params.OrderType
	if params.OrderType == "STOP_MARKET" {
		baseOrderType = "MARKET"
	} else if params.OrderType == "STOP_LIMIT" {
		baseOrderType = "LIMIT"
	}

	orderTypeSupported := false
	for _, supportedType := range contractInfo.OrderTypes {
		if supportedType == baseOrderType {
			orderTypeSupported = true
			break
		}
	}

	if !orderTypeSupported {
		return nil, fmt.Errorf("order type %s not supported for symbol %s",
			baseOrderType, params.Symbol)
	}

	// Determine default margin asset based on contract info
	marginAsset := contractInfo.QuoteAsset
	if len(contractInfo.MarginAssets) > 0 {
		marginAsset = contractInfo.MarginAssets[0]
	}

	// Set up order parameters
	orderParams := PlaceOrderParams{
		Symbol:      params.Symbol,
		Side:        params.Side,
		Type:        params.OrderType,
		Quantity:    quantity,
		MarginAsset: marginAsset, // Use the correct margin asset from contract info
		ReduceOnly:  params.ReduceOnly,
		PositionID:  params.PositionID,
	}

	// For limit orders, round the price to the correct precision
	if (params.OrderType == "LIMIT" || params.OrderType == "STOP_LIMIT") && params.Price > 0 {
		orderParams.Price = roundToDecimal(params.Price, contractInfo.PricePrecision)
	}

	// For stop orders, set the stop price
	if (params.OrderType == "STOP_MARKET" || params.OrderType == "STOP_LIMIT") && params.StopPrice > 0 {
		orderParams.StopPrice = roundToDecimal(params.StopPrice, contractInfo.PricePrecision)
	}

	log.Default().Printf("Placing order with params: %+v\n", orderParams)

	// Place the order using the standard PlaceOrder method
	responseMap, err := api.PlaceOrder(orderParams)
	if err != nil {
		return nil, err
	}

	// Convert the map response to a structured OrderResponse
	var orderResponse OrderResponse

	// Convert the map to JSON
	jsonData, err := json.Marshal(responseMap)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %v", err)
	}

	// Parse JSON into OrderResponse struct
	if err := json.Unmarshal(jsonData, &orderResponse); err != nil {
		return nil, fmt.Errorf("error parsing response into OrderResponse: %v", err)
	}

	return &orderResponse, nil
}
func (api *OrderAPI) BulletMap(params BulletParams) (OrderResponse, error) {
	// Get contract info for the symbol
	contractInfo, ok := api.client.ExchangeInfo[params.Symbol]
	if !ok {
		return OrderResponse{}, fmt.Errorf("symbol %s not found in exchange info", params.Symbol)
	}

	// Validate order type
	validOrderTypes := []OrderType{"MARKET", "LIMIT", "STOP_MARKET", "STOP_LIMIT"}
	isValidType := false
	for _, orderType := range validOrderTypes {
		if params.OrderType == orderType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return OrderResponse{}, fmt.Errorf("invalid order type: %s. Must be MARKET, LIMIT, STOP_MARKET, or STOP_LIMIT", params.OrderType)
	}

	// Validate required parameters for specific order types
	if (params.OrderType == "LIMIT" || params.OrderType == "STOP_LIMIT") && params.Price <= 0 {
		return OrderResponse{}, fmt.Errorf("price must be specified and greater than 0 for %s orders", params.OrderType)
	}

	if (params.OrderType == "STOP_MARKET" || params.OrderType == "STOP_LIMIT") && params.StopPrice <= 0 {
		return OrderResponse{}, fmt.Errorf("stopPrice must be specified and greater than 0 for %s orders", params.OrderType)
	}

	// Determine the minimum quantity based on order type
	var minQuantity float64
	var maxQuantity float64

	if params.OrderType == "MARKET" || params.OrderType == "STOP_MARKET" {
		minQuantity = contractInfo.MarketMinQuantity
		maxQuantity = contractInfo.MarketMaxQuantity
	} else {
		minQuantity = contractInfo.MinQuantity
		maxQuantity = contractInfo.MaxQuantity
	}

	// If min quantity is not set (could happen if filter parsing failed), use a safe default
	if minQuantity <= 0 {
		log.Default().Printf("Warning: Minimum quantity for %s not set, using default value\n", params.Symbol)
		minQuantity = 0.001 // Default fallback
	}

	// Calculate quantity based on minimum quantity and count
	quantity := minQuantity * params.Count

	// Check if quantity exceeds the maximum
	if maxQuantity > 0 && quantity > maxQuantity {
		return OrderResponse{}, fmt.Errorf("calculated quantity %.8f exceeds maximum allowed %.8f for %s",
			quantity, maxQuantity, params.Symbol)
	}

	// Round to the correct precision
	quantity = roundToDecimal(quantity, contractInfo.QuantityPrecision)

	// Check if the order type is supported for this symbol
	// For stop orders, we check if the base type (MARKET/LIMIT) is supported
	baseOrderType := params.OrderType
	if params.OrderType == "STOP_MARKET" {
		baseOrderType = "MARKET"
	} else if params.OrderType == "STOP_LIMIT" {
		baseOrderType = "LIMIT"
	}

	orderTypeSupported := false
	for _, supportedType := range contractInfo.OrderTypes {
		if supportedType == baseOrderType {
			orderTypeSupported = true
			break
		}
	}

	if !orderTypeSupported {
		return OrderResponse{}, fmt.Errorf("order type %s not supported for symbol %s",
			baseOrderType, params.Symbol)
	}

	// Determine default margin asset based on contract info
	marginAsset := contractInfo.QuoteAsset
	if len(contractInfo.MarginAssets) > 0 {
		marginAsset = contractInfo.MarginAssets[0]
	}

	// Set up order parameters
	orderParams := PlaceOrderParams{
		Symbol:      params.Symbol,
		Side:        params.Side,
		Type:        params.OrderType,
		Quantity:    quantity,
		PlaceType:   "ORDER_FORM",
		MarginAsset: marginAsset, // Use the correct margin asset from contract info
		ReduceOnly:  params.ReduceOnly,
		Leverage:    params.Leverage,
	}

	// For limit orders, round the price to the correct precision
	if (params.OrderType == "LIMIT" || params.OrderType == "STOP_LIMIT") && params.Price > 0 {
		orderParams.Price = roundToDecimal(params.Price, contractInfo.PricePrecision)
	}

	// For stop orders, set the stop price
	if (params.OrderType == "STOP_MARKET" || params.OrderType == "STOP_LIMIT") && params.StopPrice > 0 {
		orderParams.StopPrice = roundToDecimal(params.StopPrice, contractInfo.PricePrecision)
	}

	log.Default().Printf("Placing order with params: %+v\n", orderParams)

	// Place the order using the standard PlaceOrder method
	return api.PlaceOrder(orderParams)
}

// roundToDecimal rounds a float to the specified decimal places
func roundToDecimal(value float64, precision int) float64 {
	multiplier := math.Pow10(precision)
	return math.Round(value*multiplier) / multiplier
}
