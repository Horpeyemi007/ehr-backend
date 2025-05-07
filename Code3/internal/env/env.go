package env

import (
	"os"
	"strconv"
)

func GetString(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func GetInt(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	intValue, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func GetPrivateKeyPath() string {
	return GetString("PRIVATE_KEY_PATH", "internal/keys/private.pem")
}

func GetPublicKeyPath() string {
	return GetString("PUBLIC_KEY_PATH", "internal/keys/public.pem")
}
