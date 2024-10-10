package utils

import "go/format"

// FormatAsGoCode formats the given string as a go source file.
func FormatAsGoCode(src string) (string, error) {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}
