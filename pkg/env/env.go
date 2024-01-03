package env

import "os"

func Get(key string) string {
	return os.Getenv(key)
}

func Lookup(key string) (string, bool) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		ok = false
	}
	return value, ok
}
