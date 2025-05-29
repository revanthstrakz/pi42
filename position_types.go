package pi42

import "time"

// PositionStatus represents position status
type PositionStatus string

// Common position status values
const (
	PositionStatusOpen   PositionStatus = "OPEN"
	PositionStatusClosed PositionStatus = "CLOSED"
)

// PositionResponse represents a trading position
type PositionResponse struct {
	ID                          int      `json:"id"`
	ContractPair                string   `json:"contractPair"`
	ContractType                string   `json:"contractType"`
	EntryPrice                  float64  `json:"entryPrice"`
	Leverage                    int      `json:"leverage"`
	LiquidationPrice            float64  `json:"liquidationPrice"`
	MarginType                  string   `json:"marginType"`
	Margin                      float64  `json:"margin"`
	MarginInMarginAsset         float64  `json:"marginInMarginAsset"`
	PositionAmount              float64  `json:"positionAmount"`
	PositionID                  string   `json:"positionId"`
	PositionSize                float64  `json:"positionSize"`
	PositionStatus              string   `json:"positionStatus"`
	PositionType                string   `json:"positionType"`
	RealizedProfit              *float64 `json:"realizedProfit"` // Nullable
	Quantity                    float64  `json:"quantity"`
	BaseAsset                   string   `json:"baseAsset"`
	MarginAsset                 string   `json:"marginAsset"`
	QuoteAsset                  string   `json:"quoteAsset"`
	CreatedTime                 string   `json:"createdTime,omitempty"`
	UpdatedTime                 string   `json:"updatedTime,omitempty"`
	CreatedAt                   string   `json:"createdAt,omitempty"`
	IconUrl                     string   `json:"iconUrl,omitempty"`
	MaintenanceMarginPercentage *float64 `json:"maintenanceMarginPercentage,omitempty"`
	MarginConversionRate        *float64 `json:"marginConversionRate,omitempty"`
	MarginSettlementRate        *float64 `json:"marginSettlementRate,omitempty"`
	RealizedProfitInMarginAsset *float64 `json:"realizedProfitInMarginAsset,omitempty"`
}

// PositionCloseResponse represents the response when closing positions
type PositionCloseResponse struct {
	Success bool                  `json:"success"`
	Data    []PositionCloseStatus `json:"data"`
}

// PositionCloseStatus represents the status of a closed position
type PositionCloseStatus struct {
	PositionID string `json:"positionId"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

// ParsedCreatedTime parses the CreatedTime field string into a time.Time object
func (p PositionResponse) ParsedCreatedTime() (time.Time, error) {
	if p.CreatedTime != "" {
		return time.Parse(time.RFC3339, p.CreatedTime)
	}
	return time.Parse(time.RFC3339, p.CreatedAt)
}

// ParsedUpdatedTime parses the UpdatedTime field string into a time.Time object
func (p PositionResponse) ParsedUpdatedTime() (time.Time, error) {
	if p.UpdatedTime != "" {
		return time.Parse(time.RFC3339, p.UpdatedTime)
	}
	return time.Parse(time.RFC3339, p.CreatedAt) // If UpdatedTime not available, use CreatedAt
}
