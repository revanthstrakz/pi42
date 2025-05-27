package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/revanthstrakz/pi42"
)

// GetOpenPositions handles request to get all open positions
func (h *Pi42Handler) GetOpenPositions(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	// Parse query parameters
	symbol := c.Query("symbol")
	pageSizeStr := c.Query("pageSize", "50")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	sortOrder := c.Query("sortOrder", "DESC")

	// Get open positions
	positions, err := h.client.Position.GetPositions("OPEN", pi42.PositionQueryParams{
		Symbol:    symbol,
		PageSize:  pageSize,
		SortOrder: sortOrder,
	})
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(positions)
}

// GetClosedPositions handles request to get closed positions
func (h *Pi42Handler) GetClosedPositions(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	// Parse query parameters
	symbol := c.Query("symbol")
	pageSizeStr := c.Query("pageSize", "50")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	sortOrder := c.Query("sortOrder", "DESC")
	startTimeStr := c.Query("startTime", "0")
	endTimeStr := c.Query("endTime", "0")
	startTime, _ := strconv.ParseInt(startTimeStr, 10, 64)
	endTime, _ := strconv.ParseInt(endTimeStr, 10, 64)

	// Get closed positions
	positions, err := h.client.Position.GetPositions("CLOSED", pi42.PositionQueryParams{
		Symbol:         symbol,
		PageSize:       pageSize,
		SortOrder:      sortOrder,
		StartTimestamp: startTime,
		EndTimestamp:   endTime,
	})
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(positions)
}

// GetPosition handles request to get a specific position
func (h *Pi42Handler) GetPosition(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	positionID := c.Params("positionId")
	if positionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Position ID is required",
		})
	}

	position, err := h.client.Position.GetPosition(positionID)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(position)
}

// CloseAllPositions handles request to close all open positions
func (h *Pi42Handler) CloseAllPositions(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	result, err := h.client.Position.CloseAllPositions()
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(result)
}
