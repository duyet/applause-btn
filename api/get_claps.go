package api

import (
	"fmt"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
)

// GetClaps api get claps for a URL
func GetClaps(c *fiber.Ctx) error {
	if c.Get("Referer", "") == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No Referer header set",
		})
	}

	sourceURL, err := utils.GetSourceURL(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or missing URL",
		})
	}

	if !utils.IsURL(sourceURL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Referer is not a valid URL: %s", sourceURL),
		})
	}

	item, err := utils.GetItem(sourceURL)
	if err != nil {
		// URL not found in database, return 0 claps
		return c.JSON(0)
	}

	return c.JSON(item.Claps)
}
