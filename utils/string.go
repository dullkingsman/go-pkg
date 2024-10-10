package utils

import (
	"strings"
	"unicode"
)

// StripPrefix removes the given prefix from the given string if it exists.
func StripPrefix(string string, prefix string) string {
	if len(string) > len(prefix) && string[:len(prefix)] == prefix {
		return string[len(prefix):]
	}

	return string
}

// StripSuffix removes the given suffix from the given string if it exists.
func StripSuffix(string string, suffix string) string {
	if len(string) > len(suffix) && string[len(string)-len(suffix):] == suffix {
		return string[:len(string)-len(suffix)]
	}

	return string
}

// GetEquivalentWhiteSpace returns a string with only whitespace with the same
// length as the given string
func GetEquivalentWhiteSpace(string string) string {
	var tmp = ""

	for range string {
		tmp = tmp + " "
	}

	return tmp
}

// jsonNode is a struct that represents a node in a json string
type jsonNode struct {
	value    string
	nodeType string // string, boolean, null, number, structure, field
}

// colorJsonNode colors the given json node based on its type
func colorJsonNode(node *jsonNode) {
	switch node.nodeType {
	case "null":
		node.value = MagentaString(node.value)
		break
	case "boolean":
		node.value = YellowString(node.value)
		break
	case "number":
		node.value = LightBlueString(node.value)
		break
	case "string":
		node.value = GreenString(node.value)
		break
	}
}

// BeautifyIndentedJson colors an indented json using MarshalIndented
// and also trims the extra new lines put on every new line.
//
// ___
//
// Node: This function only supports prefix and indentation
// strings that contain only spaces and tabs.
func BeautifyIndentedJson(indentedJson []byte) string {
	var lines = strings.Split(string(indentedJson), "\n\n")

	var accumulating = "structure"

	for index, line := range lines {
		var nodes = []jsonNode{}

		for idx, char := range line {
			switch char {
			case '{', '}':
				accumulating = "structure"

				if idx == 0 {
					nodes = append(nodes, jsonNode{
						nodeType: "structure",
						value:    string(char),
					})
				} else {
					if accumulating == "string" {
						nodes[len(nodes)-1].value = nodes[len(nodes)-1].value + string(char)
					} else {
						nodes = append(nodes, jsonNode{
							nodeType: "structure",
							value:    string(char),
						})
					}
				}
			case '"':
				if accumulating == "string" || accumulating == "field" {
					if idx == len(line)-1 || idx == len(line)-2 {
						if idx == len(line)-1 {
							accumulating = "structure"
						} else {
							accumulating = "comma"
						}
					} else if line[idx-1] != '\\' {
						accumulating = "colon"
					}

					nodes[len(nodes)-1].value = nodes[len(nodes)-1].value + string(char)
				} else {
					var nodeType = "field"

					if accumulating == "colon" || accumulating == "comma" || accumulating == "array_open" {
						nodeType = "string"
					}

					accumulating = nodeType

					nodes = append(nodes, jsonNode{
						nodeType: nodeType,
						value:    string(char),
					})
				}
			case ':', ',', '[', ']', ' ', '\t', '\n', '\r':
				if accumulating == "string" {
					nodes[len(nodes)-1].value = nodes[len(nodes)-1].value + string(char)
				} else {
					if char == '[' {
						accumulating = "array_open"
					} else if char == ']' || char == ',' {
						accumulating = "structure"
					}

					nodes = append(nodes, jsonNode{
						nodeType: "structure",
						value:    string(char),
					})
				}
			default:
				if accumulating == "structure" || accumulating == "array_open" || accumulating == "colon" || accumulating == "comma" {
					var nodeType = "structure"

					if char == 'n' {
						nodeType = "null"
					} else if char == 'f' || char == 't' {
						nodeType = "boolean"
					} else if unicode.IsDigit(char) {
						nodeType = "number"
					}

					accumulating = nodeType

					nodes = append(nodes, jsonNode{
						nodeType: nodeType,
						value:    string(char),
					})
				} else {
					nodes[len(nodes)-1].value = nodes[len(nodes)-1].value + string(char)
				}
			}
		}

		var acc = ""

		for _, node := range nodes {
			colorJsonNode(&node)
			acc = acc + node.value
		}

		lines[index] = acc
	}

	return strings.Join(lines, "\n")
}

// StringArrayToCommaSeparatedString converts a string array to a comma separated string
func StringArrayToCommaSeparatedString(array []string) string {
	return strings.Join(array, ", ")
}

// SnakeCaseToPascalCase converts a snake case string to a pascal case string
func SnakeCaseToPascalCase(snake string) string {
	parts := strings.Split(snake, "_")

	for i, part := range parts {
		parts[i] = strings.Title(part)
	}

	return strings.Join(parts, "")
}

// SnakeCaseToCamelCase converts a snake case string to a camel case string
func SnakeCaseToCamelCase(s string) string {
	parts := strings.Split(s, "_")

	for i, part := range parts {
		if i > 0 {
			parts[i] = strings.Title(part)
		}
	}

	return strings.Join(parts, "")
}

// CapitalizeFirstLetter capitalizes the first letter of the given string
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// LowercaseFirstLetter lowercases the first letter of the given string
func LowercaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// CapitalizeFirstLetterOnly capitalizes the first letter of the given string
// and ensures all others are lowercase
func CapitalizeFirstLetterOnly(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// RemoveAnySpecialCharactersAndReturnSpaceSeparatedString removes any special characters
// from the string and replaces them with spaces in the returned rune slice.
//
// ___
//
// Note: Special characters are non-alphanumeric characters.
func RemoveAnySpecialCharactersAndReturnSpaceSeparatedString(s string) []rune {
	var result []rune

	if len(s) == 0 {
		return result
	}

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			result = append(result, r)
		} else {
			result = append(result, ' ')
		}
	}

	return result
}

// AnyToCamelCase converts any string to a camel case string
func AnyToCamelCase(s string) string {
	var result []rune

	if len(s) == 0 {
		return ""
	}

	var parts = strings.Split(string(RemoveAnySpecialCharactersAndReturnSpaceSeparatedString(s)), " ")

	for i, part := range parts {
		if i == 0 {
			result = append(result, []rune(strings.ToLower(part))...)
		} else {
			result = append(result, []rune(CapitalizeFirstLetterOnly(part))...)
		}
	}

	return string(result)
}

// AnyToPascalCase converts any string to a pascal case string
func AnyToPascalCase(s string) string {
	var result []rune

	if len(s) == 0 {
		return ""
	}

	var parts = strings.Split(string(RemoveAnySpecialCharactersAndReturnSpaceSeparatedString(s)), " ")

	for _, part := range parts {
		result = append(result, []rune(CapitalizeFirstLetterOnly(part))...)
	}

	return string(result)
}
