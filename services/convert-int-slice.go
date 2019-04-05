package services

import (
	"errors"
	"strconv"
	"strings"
)

// ConvertIntSlice converts []string to []int
func ConvertIntSlice(toConvert []string) ([]int, error) {
	var converted []int
	var sb strings.Builder

	for _, s := range toConvert {
		intValue, err := strconv.Atoi(s)
		if err != nil {
			sb.WriteString("Unable to convert \"")
			sb.WriteString(s)
			sb.WriteString("\" to a number")
			return nil, errors.New(sb.String())
		}
		converted = append(converted, intValue)
	}

	return converted, nil
}
