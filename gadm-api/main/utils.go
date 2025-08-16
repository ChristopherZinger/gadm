package main

import "fmt"

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func arrayToStrings(args ...interface{}) []string {
	var result []string
	for _, item := range args {
		result = append(result, fmt.Sprintf("%s", item))
	}
	return result
}
