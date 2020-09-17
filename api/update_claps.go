package api

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
)

// UpdateClaps update claps
func UpdateClaps(c *fiber.Ctx) error {
	sourceURL, err := utils.GetSourceURL(c)
	if err != nil {
		return err
	}

	body := c.Body()
	claps, err := strconv.Atoi(strings.Split(string(body), ",")[0])
	// for the v2.0.0 behavior, where the clap count was a temporal offset, always
	// treat this as a single clap
	if err != nil {
		claps = 1
	}

	if !utils.IsURL(sourceURL) {
		return fmt.Errorf("Referer is not a URL [%s]", sourceURL)
	}

	clapIncrement := utils.Clamp(claps, 1, 10)
	// var totalClaps int

	log.Printf("Adding %v claps to %s", clapIncrement, sourceURL)
	sourceIP := c.IP()

	item, err := utils.GetItem(sourceURL)
	if err != nil {
		newItem := utils.Item{SourceIP: sourceIP, Claps: clapIncrement}
		if err := utils.PutItem(sourceURL, newItem); err != nil {
			return err
		}
	} else {
		if item.SourceIP != "" && item.SourceIP == sourceIP {
			return fmt.Errorf("multiple claps from the same sourceIp prohibited %s", sourceIP)
		}

		item.Claps += clapIncrement
		item.SourceIP = sourceIP
		if err := utils.PutItem(sourceURL, item); err != nil {
			return err
		}
	}

	fmt.Printf("%s   %v  %s", sourceURL, item, sourceIP)
	return c.JSON(item.Claps)
}
