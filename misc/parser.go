// Parse printf format string tokens
package misc

import "strings"

const symbols = "vT%tbcdoqxXUeEfFgGsqxX"
const consymbols = "0123456789.+-#"

func ParseFormat(format string) (ret []string) {
	ret = make([]string, 0)

	cur := ""
	in := false
	for _, c := range format {
		if !in && (c == '%') {
			in = true
		} else if in && strings.Contains(symbols, string(c)) {
			cur += string(c)
			in = false
			ret = append(ret, cur)
			cur = ""
		} else if in && strings.Contains(consymbols, string(c)) {
			cur += string(c)
		}
	}

	return
}
