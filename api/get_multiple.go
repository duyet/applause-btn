package api

import (
	"encoding/json"
	"errors"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
)

// GetMultiple get multiple url
func GetMultiple(c *fiber.Ctx) error {
	body := c.Body()

	var listURL []string
	err := json.Unmarshal([]byte(body), &listURL)
	if err != nil {
		return errors.New("getMultiple requires an array")
	}

	// check every item is url
	for _, x := range listURL {
		if utils.IsURL(x) == false {
			return errors.New("getMultiple requires an array of URLs")
		}
	}

	if len(listURL) == 0 {
		return c.JSON([]string{})
	}

	// TODO: limit the query to 100 URLs
	item, err := utils.GetItems(listURL)
	if err != nil {
		return errors.New("Cannot get items")
	}

	return c.JSON(item)
}
