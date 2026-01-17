package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	return Config{
		Port: getenv("PORT", "8080"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
