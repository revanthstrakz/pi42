package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revanthstrakz/pi42"
)

// GetTicker handles request for ticker data for a symbol
func (h *Pi42Handler) GetTicker(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	if symbol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Symbol is required",
		})
	}

	ticker, err := h.client.Market.GetTicker24hr(symbol)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(ticker)
}

// GetDepth handles request for order book depth for a symbol
func (h *Pi42Handler) GetDepth(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	if symbol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Symbol is required",
		})
	}

	depth, err := h.client.Market.GetDepth(symbol)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(depth)
}

// GetKlines handles request for candlestick data
func (h *Pi42Handler) GetKlines(c *fiber.Ctx) error {
	// Parse request body
	params := new(struct {
		Pair      string `json:"pair"`
		Interval  string `json:"interval"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
		Limit     int    `json:"limit"`
	})

	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if params.Pair == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Pair is required",
		})
	}

	if params.Interval == "" {
		params.Interval = "1h" // Default interval
	}

	klines, err := h.client.Market.GetKlines(pi42.KlinesParams{
		Pair:      params.Pair,
		Interval:  params.Interval,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
		Limit:     params.Limit,
	})
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(klines)
}
