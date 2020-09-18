package api

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/duyet/applause-btn/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

	log.Printf("Adding %v claps to %s", clapIncrement, sourceURL)
	sourceIP := c.IP()
	clapperInfo := getClapperInfo(c)

	item, err := utils.GetItem(sourceURL)
	if err != nil {
		clappers := []utils.ClapperInfo{}
		if clapperInfo != nil {
			clappers = []utils.ClapperInfo{*clapperInfo}
		}

		newItem := utils.Item{
			SourceIP: sourceIP,
			Claps:    clapIncrement,
			Clappers: clappers,
		}
		if err := utils.PutItem(sourceURL, newItem); err != nil {
			return err
		}
	} else {
		if item.SourceIP != "" && item.SourceIP == sourceIP {
			return fmt.Errorf("multiple claps from the same sourceIp prohibited %s", sourceIP)
		}

		item.Claps += clapIncrement
		item.SourceIP = sourceIP
		item.Clappers = appendToList(item.Clappers, clapperInfo)
		if err := utils.PutItem(sourceURL, item); err != nil {
			return err
		}
	}

	fmt.Printf("%s   %v  %s", sourceURL, item, sourceIP)
	return c.JSON(item.Claps)
}

func getClapperInfo(c *fiber.Ctx) *utils.ClapperInfo {
	headerUserEmail := getEnv("HEADER_USER_EMAIL", "x-authenticated-user-email")
	headerUserID := getEnv("HEADER_USER_ID", "x-authenticated-uid")

	email := c.Get(headerUserEmail)
	uid := c.Get(headerUserID)

	if email != "" || uid != "" {
		return &utils.ClapperInfo{Email: email, UID: uid, CreatedAt: time.Now()}
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func appendToList(source []utils.ClapperInfo, new *utils.ClapperInfo) []utils.ClapperInfo {
	if new == nil {
		return source
	}

	for _, el := range source {
		if cmp.Equal(el, new, cmpopts.IgnoreFields(utils.ClapperInfo{}, "CreatedAt")) {
			return source
		}
	}
	return append(source, *new)
}
