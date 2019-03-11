package services

import "strings"

// TrimSlice Trim spaces from []string
func TrimSlice(toTrim *[]string) {
	// remove spaces from elements
	for i := range *toTrim {
		(*toTrim)[i] = strings.TrimSpace((*toTrim)[i])
	}
}
