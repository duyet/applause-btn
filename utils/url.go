package utils

import (
	"errors"
	"log"
	"net/url"

	"github.com/PuerkitoBio/purell"
	"github.com/gofiber/fiber"
)

// IsURL check is url or not
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Host != ""
}

// GetSourceURL set source url from query param or Referer header
func GetSourceURL(c *fiber.Ctx) (sourceURL string, err error) {
	sourceURL = c.Query("url", c.Get("Referer"))

	if sourceURL == "" {
		err = errors.New("no referer or url specified")
		return
	}

	normalized, err := purell.NormalizeURLString(
		sourceURL,
		purell.FlagLowercaseScheme|
			purell.FlagLowercaseHost|
			purell.FlagUppercaseEscapes|
			purell.FlagRemoveDuplicateSlashes|
			purell.FlagRemoveFragment|
			purell.FlagsUsuallySafeGreedy)
	if err != nil {
		err = errors.New("no referer or url specified")
		return
	}

	log.Printf("source URL: '%s', normalized source URL: '%s'", sourceURL, normalized)

	sourceURL = normalized

	return
}
