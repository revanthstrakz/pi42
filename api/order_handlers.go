package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/revanthstrakz/pi42"
)

// PlaceOrder handles request to place a new order
func (h *Pi42Handler) PlaceOrder(c *fiber.Ctx) error {
	// Extract API key and secret from headers (or JWT claims)
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")

	// If API credentials are provided, update the client
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	// Parse request body for order parameters
	orderReq := new(struct {
		Symbol     string  `json:"symbol"`
		Side       string  `json:"side"`
		OrderType  string  `json:"orderType"`
		Price      float64 `json:"price"`
		StopPrice  float64 `json:"stopPrice"`
		Count      float64 `json:"count"`
		ReduceOnly bool    `json:"reduceOnly"`
		PositionID string  `json:"positionId"`
	})

	if err := c.BodyParser(orderReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if orderReq.Symbol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Symbol is required",
		})
	}

	if orderReq.Side == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Side is required",
		})
	}

	if orderReq.OrderType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OrderType is required",
		})
	}

	if orderReq.Count <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Count must be greater than 0",
		})
	}

	// Use the Bullet API to place the order
	order, err := h.client.Order.Bullet(pi42.BulletParams{
		Symbol:     orderReq.Symbol,
		Side:       orderReq.Side,
		OrderType:  orderReq.OrderType,
		Price:      orderReq.Price,
		StopPrice:  orderReq.StopPrice,
		Count:      orderReq.Count,
		ReduceOnly: orderReq.ReduceOnly,
		PositionID: orderReq.PositionID,
	})
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(order)
}

// GetOpenOrders handles request to get all open orders
func (h *Pi42Handler) GetOpenOrders(c *fiber.Ctx) error {
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

	// Get open orders
	orders, err := h.client.Order.GetOpenOrders(pi42.OrderQueryParams{
		Symbol:    symbol,
		PageSize:  pageSize,
		SortOrder: sortOrder,
	})
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(orders)
}

// GetOrderHistory handles request to get order history
func (h *Pi42Handler) GetOrderHistory(c *fiber.Ctx) error {
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

	// Get order history
	orders, err := h.client.Order.GetOrderHistory(pi42.OrderQueryParams{
		Symbol:         symbol,
		PageSize:       pageSize,
		SortOrder:      sortOrder,
		StartTimestamp: startTime,
		EndTimestamp:   endTime,
	})
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(orders)
}

// CancelOrder handles request to cancel a specific order
func (h *Pi42Handler) CancelOrder(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	clientOrderID := c.Params("clientOrderId")
	if clientOrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Client order ID is required",
		})
	}

	result, err := h.client.Order.DeleteOrder(clientOrderID)
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(result)
}

// CancelAllOrders handles request to cancel all open orders
func (h *Pi42Handler) CancelAllOrders(c *fiber.Ctx) error {
	// Extract API credentials if provided
	apiKey := c.Get("X-API-Key")
	apiSecret := c.Get("X-API-Secret")
	if apiKey != "" && apiSecret != "" {
		h.SetClient(apiKey, apiSecret)
	}

	result, err := h.client.Order.CancelAllOrders()
	if err != nil {
		return handleApiError(c, err)
	}

	return c.JSON(result)
}
