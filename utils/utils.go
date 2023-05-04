package utils

import (
	"os"
	"strconv"
)

func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func GetenvInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)

	if err != nil {
		return defaultValue
	}

	return intValue
}
