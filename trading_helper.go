package pi42

import (
	"fmt"
	"math"
	"strconv"
)

// TradingHelper provides convenient access to symbol-specific trading parameters
// like minimum quantities, price precision, and price step sizes
type TradingHelper struct {
	// Symbol information
	Symbol       string
	BaseAsset    string
	QuoteAsset   string
	MarginAsset  string
	ContractType string

	// Quantity constraints
	MinQuantity       float64
	MaxQuantity       float64
	QuantityPrecision int

	// Price constraints
	MinPrice       float64 // Minimum valid price (often 0)
	MaxPrice       float64 // Maximum valid price (often very high)
	PricePrecision int     // Number of decimal places for price
	MinPriceStep   float64 // Minimum price increment

	// Derived values
	PercentIncrement float64 // Percentage of price difference between steps

	// Reference to client for market data access
	client *Client
}

// NewTradingHelper creates a new TradingHelper for a specific symbol
// The percentIncrement parameter defines the granularity of price steps as a percentage
func NewTradingHelper(client *Client, symbol string, percentIncrement float64) (*TradingHelper, error) {
	// Create new helper instance
	helper := &TradingHelper{
		Symbol:           symbol,
		PercentIncrement: percentIncrement,
		client:           client,
	}

	// Initialize the helper with contract specifications
	if err := helper.init(); err != nil {
		return nil, err
	}

	return helper, nil
}

// init loads all necessary trading parameters from the exchange
func (th *TradingHelper) init() error {
	// First check if we already have the contract info cached in the client
	contractInfo, exists := th.client.ExchangeInfo[th.Symbol]
	if !exists {
		// If not cached, try to fetch exchange info
		if err := th.client.fetchExchangeInfo(); err != nil {
			return fmt.Errorf("failed to fetch exchange info: %v", err)
		}

		// Check again after fetching
		contractInfo, exists = th.client.ExchangeInfo[th.Symbol]
		if !exists {
			return fmt.Errorf("symbol %s not found in exchange info", th.Symbol)
		}
	}

	// Populate fields from contract info
	th.BaseAsset = contractInfo.BaseAsset
	th.QuoteAsset = contractInfo.QuoteAsset
	th.ContractType = contractInfo.ContractType
	th.QuantityPrecision = contractInfo.QuantityPrecision
	th.PricePrecision = contractInfo.PricePrecision
	th.MinQuantity = contractInfo.MinQuantity
	th.MaxQuantity = contractInfo.MaxQuantity

	// Set default margin asset if available
	if len(contractInfo.MarginAssets) > 0 {
		th.MarginAsset = contractInfo.MarginAssets[0]
	} else {
		th.MarginAsset = contractInfo.QuoteAsset
	}

	// Calculate minimum price step based on precision
	th.MinPriceStep = 1.0 / math.Pow10(th.PricePrecision)

	// Get current market price to calculate percentage-based increments
	if err := th.updateCurrentPrice(); err != nil {
		return fmt.Errorf("failed to get current price: %v", err)
	}

	return nil
}

// updateCurrentPrice gets the latest market price for the symbol
func (th *TradingHelper) updateCurrentPrice() error {
	ticker, err := th.client.Market.GetTicker24hr(th.Symbol)
	if err != nil {
		return err
	}

	data, ok := ticker["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("could not parse ticker data")
	}

	lastPrice, ok := data["c"].(string)
	if !ok {
		return fmt.Errorf("could not parse last price")
	}

	currentPrice, err := strconv.ParseFloat(lastPrice, 64)
	if err != nil {
		return fmt.Errorf("could not convert price to float: %v", err)
	}

	// Set MaxPrice to a high multiple of current price
	th.MaxPrice = currentPrice * 10

	return nil
}

// GetMinimumOrderQuantity returns the minimum quantity allowed for orders
func (th *TradingHelper) GetMinimumOrderQuantity() float64 {
	return th.MinQuantity
}

// GetMinimumPriceIncrement returns the smallest price step allowed
func (th *TradingHelper) GetMinimumPriceIncrement() float64 {
	return th.MinPriceStep
}

// GetPricePrecision returns the number of decimal places for price
func (th *TradingHelper) GetPricePrecision() int {
	return th.PricePrecision
}

// GetQuantityPrecision returns the number of decimal places for quantity
func (th *TradingHelper) GetQuantityPrecision() int {
	return th.QuantityPrecision
}

// GetMarginAsset returns the default margin asset for this symbol
func (th *TradingHelper) GetMarginAsset() string {
	return th.MarginAsset

}

// SymbolInfo contains basic information about a trading symbol
type SymbolInfo struct {
	Symbol            string  `json:"symbol"`
	BaseAsset         string  `json:"baseAsset"`
	QuoteAsset        string  `json:"quoteAsset"`
	MarginAsset       string  `json:"marginAsset"`
	ContractType      string  `json:"contractType"`
	MinQuantity       float64 `json:"minQuantity"`
	MaxQuantity       float64 `json:"maxQuantity"`
	QuantityPrecision int     `json:"quantityPrecision"`
	PricePrecision    int     `json:"pricePrecision"`
	MinPriceStep      float64 `json:"minPriceStep"`
	PercentIncrement  float64 `json:"percentIncrement"`
}

// GetSymbolInfo returns basic information about the symbol
func (th *TradingHelper) GetSymbolInfo() SymbolInfo {
	return SymbolInfo{
		Symbol:            th.Symbol,
		BaseAsset:         th.BaseAsset,
		QuoteAsset:        th.QuoteAsset,
		MarginAsset:       th.MarginAsset,
		ContractType:      th.ContractType,
		MinQuantity:       th.MinQuantity,
		MaxQuantity:       th.MaxQuantity,
		QuantityPrecision: th.QuantityPrecision,
		PricePrecision:    th.PricePrecision,
		MinPriceStep:      th.MinPriceStep,
		PercentIncrement:  th.PercentIncrement,
	}
}

// CalculatePriceFromBestPrice calculates a price at a specified percentage difference
// from the best bid/ask price. Positive percentDiff for above, negative for below.
func (th *TradingHelper) CalculatePriceFromBestPrice(percentDiff float64) (float64, error) {
	// Get depth data to find best bid/ask
	depth, err := th.client.Market.GetDepth(th.Symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get order book depth: %v", err)
	}

	var bestPrice float64

	// For positive percentDiff, we start from the best ask (for buy orders)
	// For negative percentDiff, we start from the best bid (for sell orders)
	if percentDiff > 0 {
		// Use best ask (lowest sell price) as reference
		if len(depth.Data.Asks) > 0 {
			bestPrice, err = strconv.ParseFloat(depth.Data.Asks[0][0], 64)
			if err != nil {
				return 0, fmt.Errorf("could not parse ask price: %v", err)
			}
		} else {
			return 0, fmt.Errorf("no ask prices available in order book")
		}
	} else {
		// Use best bid (highest buy price) as reference
		if len(depth.Data.Bids) > 0 {
			bestPrice, err = strconv.ParseFloat(depth.Data.Bids[0][0], 64)
			if err != nil {
				return 0, fmt.Errorf("could not parse bid price: %v", err)
			}
		} else {
			return 0, fmt.Errorf("no bid prices available in order book")
		}
	}

	// Calculate target price with percentage difference
	targetPrice := bestPrice * (1 + percentDiff/100)

	// Round to the correct precision
	targetPrice = math.Round(targetPrice/th.MinPriceStep) * th.MinPriceStep

	return targetPrice, nil
}

// GetCurrentBestPrices returns the current best bid and ask prices
func (th *TradingHelper) GetCurrentBestPrices() (float64, float64, error) {
	// Get depth data to find best bid/ask
	depth, err := th.client.Market.GetDepth(th.Symbol)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get order book depth: %v", err)
	}

	var bestBid, bestAsk float64

	// Get best bid (highest buy price)
	if len(depth.Data.Bids) > 0 {
		bestBid, err = strconv.ParseFloat(depth.Data.Bids[0][0], 64)
		if err != nil {
			return 0, 0, fmt.Errorf("could not parse bid price: %v", err)
		}
	} else {
		return 0, 0, fmt.Errorf("no bid prices available in order book")
	}

	// Get best ask (lowest sell price)
	if len(depth.Data.Asks) > 0 {
		bestAsk, err = strconv.ParseFloat(depth.Data.Asks[0][0], 64)
		if err != nil {
			return 0, 0, fmt.Errorf("could not parse ask price: %v", err)
		}
	} else {
		return 0, 0, fmt.Errorf("no ask prices available in order book")
	}

	return bestBid, bestAsk, nil
}

// CalculateOrderQuantity calculates order quantity in base asset units
// from an amount in quote asset (e.g., INR amount to BTC quantity)
func (th *TradingHelper) CalculateOrderQuantity(quoteAmount float64) (float64, error) {
	// Get current price to calculate conversion
	bestBid, bestAsk, err := th.GetCurrentBestPrices()
	if err != nil {
		return 0, err
	}

	// Use average of bid/ask for calculation
	averagePrice := (bestBid + bestAsk) / 2

	// Calculate quantity
	quantity := quoteAmount / averagePrice

	// Check against minimum
	if quantity < th.MinQuantity {
		return 0, fmt.Errorf("calculated quantity %.8f is below minimum allowed %.8f",
			quantity, th.MinQuantity)
	}

	// Check against maximum
	if th.MaxQuantity > 0 && quantity > th.MaxQuantity {
		return 0, fmt.Errorf("calculated quantity %.8f is above maximum allowed %.8f",
			quantity, th.MaxQuantity)
	}

	// Round to the correct precision
	precisionMultiplier := math.Pow10(th.QuantityPrecision)
	quantity = math.Floor(quantity*precisionMultiplier) / precisionMultiplier

	return quantity, nil
}
