package api

import (
	"encoding/json"
	"errors"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber"
)

// GetMultiple get multiple url
func GetMultiple(c *fiber.Ctx) {
	body := c.Body()

	var listURL []string
	err := json.Unmarshal([]byte(body), &listURL)
	if err != nil {
		c.Next(errors.New("getMultiple requires an array"))
		return
	}

	// check every item is url
	for _, x := range listURL {
		if utils.IsURL(x) == false {
			c.Next(errors.New("getMultiple requires an array of URLs"))
			return
		}
	}

	if len(listURL) == 0 {
		c.JSON([]string{})
		return
	}

	// TODO: limit the query to 100 URLs
	item, err := utils.GetItems(listURL)
	if err != nil {
		c.Next(err)
		return
	}

	c.JSON(item)
}
