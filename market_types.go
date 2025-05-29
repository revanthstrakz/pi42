package pi42

// DepthResponse represents the full response from the GetDepth endpoint
type DepthResponse struct {
	Data DepthData `json:"data"`
}

// DepthData represents the order book depth data structure
type DepthData struct {
	EventType     string     `json:"e"`  // Event type (depthUpdate)
	EventTime     int64      `json:"E"`  // Event time in milliseconds
	TransactionTs int64      `json:"T"`  // Transaction time in milliseconds
	Symbol        string     `json:"s"`  // Trading pair symbol
	FirstUpdateID int64      `json:"U"`  // First update ID in the update
	LastUpdateID  int64      `json:"u"`  // Last update ID in the update
	PrevUpdateID  int64      `json:"pu"` // Previous update ID
	Bids          [][]string `json:"b"`  // Bid prices and quantities [price, quantity][]
	Asks          [][]string `json:"a"`  // Ask prices and quantities [price, quantity][]
}

// KlineData represents a single candlestick/kline data point
type KlineData struct {
	StartTime string `json:"startTime"` // Start time of the interval in milliseconds
	Open      string `json:"open"`      // Opening price of the interval
	High      string `json:"high"`      // Highest price during the interval
	Low       string `json:"low"`       // Lowest price during the interval
	Close     string `json:"close"`     // Closing price of the interval
	EndTime   string `json:"endTime"`   // End time of the interval in milliseconds
	Volume    string `json:"volume"`    // Trading volume during the interval
}
