package pi42

import "time"

// OrderSide represents order side (BUY or SELL)
type OrderSide string

// Common order side values
const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderType represents order type
type OrderType string

// Common order types
const (
	OrderTypeMarket     OrderType = "MARKET"
	OrderTypeLimit      OrderType = "LIMIT"
	OrderTypeStopMarket OrderType = "STOP_MARKET"
	OrderTypeStopLimit  OrderType = "STOP_LIMIT"
)

// TimeInForce represents order time in force
type TimeInForce string

// Common time in force options
const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill or Kill
	TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
)

// PositionSide represents order position side
type PositionSide string

// Common position side values
const (
	PositionSideLong  PositionSide = "LONG"
	PositionSideShort PositionSide = "SHORT"
	PositionSideBoth  PositionSide = "BOTH"
)

// OrderStatus represents the status of an order
type OrderStatus string

// Common order status values
const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// OpenOrder represents an open order
type OpenOrder struct {
	ClientOrderID   string  `json:"clientOrderId"`
	Time            string  `json:"time"`
	Symbol          string  `json:"symbol"`
	ContractType    string  `json:"contractType"`
	Type            string  `json:"type"`
	Side            string  `json:"side"`
	Price           float64 `json:"price"`
	OrderAmount     float64 `json:"orderAmount"`
	FilledAmount    float64 `json:"filledAmount"`
	LinkID          string  `json:"linkId"`
	LinkType        string  `json:"linkType"`
	SubType         string  `json:"subType"`
	PlaceType       string  `json:"placeType"`
	BaseAsset       string  `json:"baseAsset"`
	QuoteAsset      string  `json:"quoteAsset"`
	Leverage        int     `json:"leverage"`
	LockedMargin    float64 `json:"lockedMargin"`
	MarginAsset     string  `json:"marginAsset,omitempty"`
	Status          string  `json:"status,omitempty"`
	StopPrice       float64 `json:"stopPrice,omitempty"`
	ReduceOnly      bool    `json:"reduceOnly,omitempty"`
	TakeProfitPrice float64 `json:"takeProfitPrice,omitempty"`
	StopLossPrice   float64 `json:"stopLossPrice,omitempty"`
}

// OrderHistoryItem represents an item in order history
type OrderHistoryItem struct {
	ClientOrderID             string  `json:"clientOrderId"`
	UpdatedAt                 string  `json:"updatedAt"`
	Symbol                    string  `json:"symbol"`
	Type                      string  `json:"type"`
	IsIsolated                bool    `json:"isIsolated"`
	Side                      string  `json:"side"`
	Price                     string  `json:"price"`
	AvgPrice                  string  `json:"avgPrice"`
	OrigQty                   string  `json:"origQty"`
	CumQty                    string  `json:"cumQty"`
	ExecutedQty               string  `json:"executedQty"`
	ReduceOnly                bool    `json:"reduceOnly"`
	Status                    string  `json:"status"`
	Leverage                  int     `json:"leverage"`
	SubType                   string  `json:"subType"`
	StopPrice                 *string `json:"stopPrice"` // Nullable
	LockedMargin              float64 `json:"lockedMargin"`
	LockedMarginInMarginAsset float64 `json:"lockedMarginInMarginAsset"`
	MarginAsset               string  `json:"marginAsset"`
	ContractType              string  `json:"contractType"`
	IconUrl                   string  `json:"iconUrl"`
	QuoteAsset                string  `json:"quoteAsset"`
	BaseAsset                 string  `json:"baseAsset"`
	LeveragedQty              float64 `json:"leveragedQty"`
}

// LinkedOrder represents an order linked to another order
type LinkedOrder struct {
	ClientOrderID   string   `json:"clientOrderId"`
	Time            string   `json:"time"`
	Symbol          string   `json:"symbol"`
	ContractType    string   `json:"contractType"`
	Type            string   `json:"type"`
	Side            string   `json:"side"`
	Price           float64  `json:"price"`
	OrderAmount     float64  `json:"orderAmount"`
	FilledAmount    float64  `json:"filledAmount"`
	LinkID          string   `json:"linkId"`
	SubType         string   `json:"subType"`
	LinkType        string   `json:"linkType"`
	TakeProfitPrice *float64 `json:"takeProfitPrice"` // Nullable
	StopLossPrice   *float64 `json:"stopLossPrice"`   // Nullable
	Status          string   `json:"status"`
	PlaceType       string   `json:"placeType"`
	BaseAsset       string   `json:"baseAsset"`
	QuoteAsset      string   `json:"quoteAsset"`
}

// OrderCancelResponse represents the response when canceling an order
type OrderCancelResponse struct {
	ClientOrderID string `json:"clientOrderId"`
	OrderID       int    `json:"orderId"`
	Status        string `json:"status"`
	Success       bool   `json:"success"`
}

// BatchCancelResponse represents the response when canceling multiple orders
type BatchCancelResponse struct {
	Success bool                     `json:"success"`
	Data    []OrderCancelationStatus `json:"data"`
}

// OrderCancelationStatus represents the status of a canceled order
type OrderCancelationStatus struct {
	ClientOrderID string `json:"clientOrderId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

// ParsedTime parses the Time field string into a time.Time object
func (o OpenOrder) ParsedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, o.Time)
}

// ParsedTime parses the UpdatedAt field string into a time.Time object
func (o OrderHistoryItem) ParsedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, o.UpdatedAt)
}

// ParsedTime parses the Time field string into a time.Time object
func (l LinkedOrder) ParsedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, l.Time)
}
