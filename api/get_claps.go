package api

import (
	"errors"
	"fmt"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
)

// GetClaps api get claps
func GetClaps(c *fiber.Ctx) error {
	if c.Get("Referer", "") == "" {
		return errors.New("no referer set")
	}

	sourceURL, err := utils.GetSourceURL(c)
	if err != nil {
		return err
	}

	if utils.IsURL(sourceURL) == false {
		return fmt.Errorf("Referer is not a URL [%s]", sourceURL)
	}

	item, err := utils.GetItem(sourceURL)
	if err != nil {
		return c.JSON(0)
	}

	return c.JSON(item.Claps)
}
