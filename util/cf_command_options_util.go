package util

import "bytes"

//GetShortOption ...
func GetShortOption(option string) string {
	var opt bytes.Buffer
	opt.WriteString(optionPrefix)
	opt.WriteString(option)
	return opt.String()
}
