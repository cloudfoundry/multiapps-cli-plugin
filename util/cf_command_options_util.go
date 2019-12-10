package util

import (
	"bytes"
	"strings"
)

//GetShortOption ...
func GetShortOption(option string) string {
	var opt bytes.Buffer
	opt.WriteString(optionPrefix)
	opt.WriteString(option)
	return opt.String()
}

func DiscardIfEmpty(value string) *string {
	if len(value) > 0 {
		return &value
	} else {
		return nil
	}
}

func TrimAndDiscardIfEmpty(value string) *string {
	var trimmedValue = strings.TrimSpace(value)
	if len(trimmedValue) > 0 {
		return &trimmedValue
	} else {
		return nil
	}
}

func NamespaceInfoTextIfApplicable(namespace string) string {
	if DiscardIfEmpty(namespace) == nil {
		return ""
	}

	return " with namespace " + namespace
}
