package api

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber"
)

// UpdateClaps update claps
func UpdateClaps(c *fiber.Ctx) {
	sourceURL, err := utils.GetSourceURL(c)
	if err != nil {
		c.Next(err)
	}

	body := c.Body()
	claps, err := strconv.Atoi(strings.Split(body, ",")[0])
	// for the v2.0.0 behavior, where the clap count was a temporal offset, always
	// treat this as a single clap
	if err != nil {
		claps = 1
	}

	if !utils.IsURL(sourceURL) {
		c.Next(fmt.Errorf("Referer is not a URL [%s]", sourceURL))
		return
	}

	clapIncrement := utils.Clamp(claps, 1, 10)
	// var totalClaps int

	log.Printf("adding %v claps to %s", clapIncrement, sourceURL)
	sourceIP := c.IP()

	item, err := utils.GetItem(sourceURL)
	if err != nil {
		newItem := utils.Item{SourceIP: sourceIP, Claps: clapIncrement}
		if err := utils.PutItem(sourceURL, newItem); err != nil {
			c.Next(err)
			return
		}
	} else {
		if item.SourceIP != "" && item.SourceIP == sourceIP {
			c.Next(fmt.Errorf("multiple claps from the same sourceIp prohibited %s", sourceIP))
			return
		}

		item.Claps += clapIncrement
		item.SourceIP = sourceIP
		if err := utils.PutItem(sourceURL, item); err != nil {
			c.Next(err)
			return
		}
	}

	fmt.Printf("%s   %v  %s", sourceURL, item, sourceIP)
}
