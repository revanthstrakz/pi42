package api

import (
	"github.com/gofiber/fiber/v2"
)

// GetFuturesWallet handles request to get futures wallet details
func (h *Pi42Handler) GetFuturesWallet(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	asset := c.Params("asset", "INR")
	if asset == "" {
		asset = "INR" // Default to INR if not specified
	}

	wallet, err := h.client.Wallet.FuturesWalletDetails(asset)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(wallet)
}

// GetFundingWallet handles request to get funding wallet details
func (h *Pi42Handler) GetFundingWallet(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	asset := c.Params("asset", "INR")
	if asset == "" {
		asset = "INR" // Default to INR if not specified
	}

	wallet, err := h.client.Wallet.FundingWalletDetails(asset)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(wallet)
}
