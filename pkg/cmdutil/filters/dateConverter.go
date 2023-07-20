package filters

import (
	"regexp"
	"strconv"
	"time"
)

func ShortDateTimeToGoDate(shortDateTime string) time.Time {
	date := time.Now()

	if shortDateTime == "now-" || shortDateTime == "now" {
		return date
	}

	re := regexp.MustCompile(`now-(\d+)([mdhs])`)
	match := re.FindStringSubmatch(shortDateTime)

	if match != nil {
		value, _ := strconv.Atoi(match[1])
		unit := match[2]

		switch unit {
		case "s":
			date = date.Add(-time.Duration(value) * time.Second)
		case "m":
			date = date.Add(-time.Duration(value) * time.Minute)
		case "h":
			date = date.Add(-time.Duration(value) * time.Hour)
		case "d":
			date = date.AddDate(0, 0, -value)
		}
	}

	return date
}
