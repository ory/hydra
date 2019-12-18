package pkg

import (
	"regexp"
)

var (
	regx *regexp.Regexp
)

func statsdSanitizerRegex() *regexp.Regexp {
	if regx != nil {
		return regx
	}
	regx, _ := regexp.Compile("[^a-zA-Z0-9._-]+")
	return regx
}

func SanitizeForStatsd(str string, replaceWith string) string {
	regex := statsdSanitizerRegex()
	return regex.ReplaceAllString(str, replaceWith)
}
