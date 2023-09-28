package ees

import "strconv"

func Float64ToString(val interface{}) string {
	if f, ok := val.(float64); ok {
		return strconv.FormatFloat(f, 'f', 0, 64)
	}
	return "0"
}
