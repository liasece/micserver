package strings

import (
	"strings"
)

func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

func MysqlRealEscapeStringBack(value string) string {
	replace := map[string]string{"\\\\": "\\", `\'`: "'", "\\\\0": "\\0", "\\n": "\n", "\\r": "\r", `\"`: `"`, "\\Z": "\x1a"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
