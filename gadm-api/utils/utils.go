package utils

import "fmt"

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ArrayToStrings(args ...interface{}) []string {
	var result []string
	for _, item := range args {
		result = append(result, fmt.Sprintf("%s", item))
	}
	return result
}
