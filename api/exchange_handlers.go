package api

import (
	"github.com/gofiber/fiber/v2"
)

// GetExchangeInfo handles request to get exchange information
func (h *Pi42Handler) GetExchangeInfo(c *fiber.Ctx) error {
	market := c.Query("market", "")
	info, err := h.client.Exchange.ExchangeInfo(market)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(info)
}

// UpdateLeverage handles request to update leverage for a contract
func (h *Pi42Handler) UpdateLeverage(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	// Parse request body
	params := new(struct {
		Leverage     int    `json:"leverage"`
		ContractName string `json:"contractName"`
	})

	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if params.Leverage <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Leverage must be greater than 0",
		})
	}

	if params.ContractName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Contract name is required",
		})
	}

	result, err := h.client.Exchange.UpdateLeverage(params.Leverage, params.ContractName)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(result)
}

// UpdatePreference handles request to update trading preference (leverage and margin mode)
func (h *Pi42Handler) UpdatePreference(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	// Parse request body
	params := new(struct {
		Leverage     int    `json:"leverage"`
		MarginMode   string `json:"marginMode"`
		ContractName string `json:"contractName"`
	})

	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if params.Leverage <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Leverage must be greater than 0",
		})
	}

	if params.MarginMode != "ISOLATED" && params.MarginMode != "CROSS" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Margin mode must be either ISOLATED or CROSS",
		})
	}

	if params.ContractName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Contract name is required",
		})
	}

	result, err := h.client.Exchange.UpdatePreference(params.Leverage, params.MarginMode, params.ContractName)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(result)
}
