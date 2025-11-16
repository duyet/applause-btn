package api

import (
	"fmt"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
)

// GetClappers api get clappers information for a URL
func GetClappers(c *fiber.Ctx) error {
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
		// URL not found in database, return empty array
		return c.JSON([]utils.ClapperInfo{})
	}

	if item.Clappers == nil {
		return c.JSON([]utils.ClapperInfo{})
	}

	return c.JSON(item.Clappers)
}
