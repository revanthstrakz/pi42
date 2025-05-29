package pi42

import "time"

// TradeHistoryItem represents an individual trade record
type TradeHistoryItem struct {
	ID             int     `json:"id"`
	Time           string  `json:"time"`
	Symbol         string  `json:"symbol"`
	Type           string  `json:"type"`
	Side           string  `json:"side"`
	Price          float64 `json:"price"`
	Quantity       float64 `json:"quantity"`
	Role           string  `json:"role"`
	Fee            float64 `json:"fee"`
	RealizedProfit float64 `json:"realizedProfit"`
	ContractType   string  `json:"contractType"`
	ClientOrderID  string  `json:"clientOrderId"`
	BaseAsset      string  `json:"baseAsset"`
	QuoteAsset     string  `json:"quoteAsset"`
	MarginAsset    string  `json:"marginAsset"`
}

// TransactionHistoryItem represents an individual transaction record
type TransactionHistoryItem struct {
	ID           int     `json:"id"`
	Time         string  `json:"time"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	Asset        string  `json:"asset"`
	Symbol       string  `json:"symbol"`
	ContractType string  `json:"contractType"`
	BaseAsset    string  `json:"baseAsset"`
	QuoteAsset   string  `json:"quoteAsset"`
}

// ParsedTime parses the Time field string into a time.Time object
func (t TradeHistoryItem) ParsedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, t.Time)
}

// ParsedTime parses the Time field string into a time.Time object
func (t TransactionHistoryItem) ParsedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, t.Time)
}
