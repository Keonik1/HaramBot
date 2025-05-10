package parse_env

import (
	"log"
	"os"
	"strconv"
)

func getEnvWithDefault[T any](key string, defaultValue *T, parser func(string) (T, error)) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		if defaultValue == nil {
			log.Fatalf("ERROR: The environment variable '%s' is required, but not set.\n", key)
		}
		return *defaultValue
	}

	parsedValue, err := parser(value)
	if err != nil {
		log.Printf("ERROR: Error parsing the environment variable '%s': %s. The default value will be used.", key, value)
		log.Printf("Error: %s", err.Error())
		if defaultValue == nil {
			log.Fatalf("ERROR: Parsing failed, and no default value was provided. Exiting.")
		}
		return *defaultValue
	}
	return parsedValue
}

func getDefaultValuePointer[T any](defaultValue ...T) *T {
	if len(defaultValue) > 0 {
		return &defaultValue[0]
	}
	return nil
}

func GetEnvString(key string, defaultValue ...string) string {
	defaultVal := getDefaultValuePointer(defaultValue...)
	return getEnvWithDefault(key, defaultVal, func(value string) (string, error) {
		return value, nil
	})
}

func GetEnvBool(key string, defaultValue ...bool) bool {
	defaultVal := getDefaultValuePointer(defaultValue...)
	return getEnvWithDefault(key, defaultVal, strconv.ParseBool)
}

func GetEnvInt(key string, defaultValue ...int) int {
	defaultVal := getDefaultValuePointer(defaultValue...)
	return getEnvWithDefault(key, defaultVal, strconv.Atoi)
}

func GetEnvFloat(key string, defaultValue ...float64) float64 {
	defaultVal := getDefaultValuePointer(defaultValue...)
	return getEnvWithDefault(key, defaultVal, func(value string) (float64, error) {
		return strconv.ParseFloat(value, 64)
	})
}
