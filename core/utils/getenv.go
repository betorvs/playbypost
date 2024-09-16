package utils

import "os"

// GetEnv receives a key and a default values string and return a string
// it looks in environment variable for a key and returns it
func GetEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	} else {
		return value
	}
}
