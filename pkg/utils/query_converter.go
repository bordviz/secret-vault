package utils

import "strings"

func QueryConvert(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\n", " "), "\t", "")
}
