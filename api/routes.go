package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revanthstrakz/pi42"
)

// Pi42Handler handles API requests for Pi42 exchange
type Pi42Handler struct {
	client *pi42.Client
}

// NewPi42Handler creates a new Pi42Handler with the provided client
func NewPi42Handler(apiKey, apiSecret string) *Pi42Handler {
	return &Pi42Handler{
		client: pi42.NewClient(apiKey, apiSecret),
	}
}

// SetupPi42Routes configures all Pi42 API routes
func SetupPi42Routes(api fiber.Router, authMiddleware fiber.Handler) {
	handler := NewPi42Handler("", "") // Will be initialized with empty credentials

	// Group routes
	marketRoutes := api.Group("/market")
	orderRoutes := api.Group("/order", authMiddleware)       // Protected
	positionRoutes := api.Group("/position", authMiddleware) // Protected
	walletRoutes := api.Group("/wallet", authMiddleware)     // Protected
	exchangeRoutes := api.Group("/exchange")

	// Market routes (public)
	marketRoutes.Get("/ticker/:symbol", handler.GetTicker)
	marketRoutes.Get("/depth/:symbol", handler.GetDepth)
	marketRoutes.Post("/klines", handler.GetKlines)

	// Order routes (authenticated)
	orderRoutes.Post("/place", handler.PlaceOrder)
	orderRoutes.Get("/open", handler.GetOpenOrders)
	orderRoutes.Get("/history", handler.GetOrderHistory)
	orderRoutes.Delete("/cancel/:clientOrderId", handler.CancelOrder)
	orderRoutes.Delete("/cancel-all", handler.CancelAllOrders)

	// Position routes (authenticated)
	positionRoutes.Get("/open", handler.GetOpenPositions)
	positionRoutes.Get("/closed", handler.GetClosedPositions)
	positionRoutes.Get("/:positionId", handler.GetPosition)
	positionRoutes.Delete("/close-all", handler.CloseAllPositions)

	// Wallet routes (authenticated)
	walletRoutes.Get("/futures/:asset", handler.GetFuturesWallet)
	walletRoutes.Get("/funding/:asset", handler.GetFundingWallet)

	// Exchange routes (public)
	exchangeRoutes.Get("/info", handler.GetExchangeInfo)
	exchangeRoutes.Post("/leverage", authMiddleware, handler.UpdateLeverage)
	exchangeRoutes.Post("/preference", authMiddleware, handler.UpdatePreference)
}

// SetClient replaces the client with one using the provided credentials
func (h *Pi42Handler) SetClient(apiKey, apiSecret string) {
	h.client = pi42.NewClient(apiKey, apiSecret)
}
