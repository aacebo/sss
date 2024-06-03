package utils

import "os"

func GetEnv(name string, defaultValue string) string {
	v := os.Getenv(name)

	if v == "" {
		return defaultValue
	}

	return v
}
