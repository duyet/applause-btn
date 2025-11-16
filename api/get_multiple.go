package api

import (
	"encoding/json"
	"fmt"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	// MaxURLsPerRequest limits the number of URLs that can be queried at once
	MaxURLsPerRequest = 100
)

// GetMultiple get multiple url
func GetMultiple(c *fiber.Ctx) error {
	body := c.Body()

	var listURL []string
	err := json.Unmarshal([]byte(body), &listURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request body must be a JSON array of URLs",
		})
	}

	if len(listURL) == 0 {
		return c.JSON([]utils.Item{})
	}

	// Limit the number of URLs that can be queried
	if len(listURL) > MaxURLsPerRequest {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fmt.Sprintf("Too many URLs requested. Maximum is %d", MaxURLsPerRequest),
			"maximum": MaxURLsPerRequest,
			"requested": len(listURL),
		})
	}

	// Validate every item is a URL
	for i, url := range listURL {
		if !utils.IsURL(url) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid URL at index %d: %s", i, url),
				"index": i,
				"url":   url,
			})
		}
	}

	items, err := utils.GetItems(listURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve items from database",
		})
	}

	return c.JSON(items)
}
