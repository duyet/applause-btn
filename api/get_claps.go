package api

import (
	"errors"
	"fmt"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber"
)

// GetClaps api get claps
func GetClaps(c *fiber.Ctx) {
	if c.Get("Referer", "") == "" {
		c.Next(errors.New("no referer set"))
		return
	}

	sourceURL, err := utils.GetSourceURL(c)
	if err != nil {
		c.Next(err)
		return
	}

	if utils.IsURL(sourceURL) == false {
		c.Next(fmt.Errorf("Referer is not a URL [%s]", sourceURL))
		return
	}

	item, err := utils.GetItem(sourceURL)
	if err != nil {
		c.JSON(0)
		return
	}

	c.JSON(item.Claps)
}
