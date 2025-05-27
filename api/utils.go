package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revanthstrakz/pi42"
)

// handleApiError processes API errors and returns an appropriate HTTP response
func handleApiError(c *fiber.Ctx, err error) error {
	// Check for Pi42 API errors
	if apiErr, ok := err.(pi42.APIError); ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     apiErr.Message,
			"errorCode": apiErr.ErrorCode,
			"details":   apiErr.Details,
		})
	}

	// Generic error handling
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
