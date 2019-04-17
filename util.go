package model

import (
	"strings"
)

func hasPrefixIgnoringCase(s string, prefix string) bool {
	return strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(prefix))
}

func hasSuffixIgnoringCase(s string, suffix string) bool {
	return strings.HasSuffix(strings.ToUpper(s), strings.ToUpper(suffix))
}
