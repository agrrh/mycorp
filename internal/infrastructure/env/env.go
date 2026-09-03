package env

import "os"

var defaultPrefix = "MYCORP_"

func GetWithPrefix(key string) string {
	prefix := os.Getenv("MYCORP_ENV_PREFIX")
	if prefix == "" {
		prefix = defaultPrefix
	}
	return os.Getenv(prefix + key)
}
