package config

type Config struct {
	ServerPort string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		ServerPort: ":8080",
		JWTSecret:  "change-me",
	}
}
