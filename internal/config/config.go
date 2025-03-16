package config

import "os"

type Config struct {
	HostnamePort string

	ApiKey     string
	SampleRate int
	BitDepth   int
	Channels   int

	PrivateKey string
	PublicKey  string

	DBPath string
}

func LoadConfig() *Config {
	return &Config{
		//server
		HostnamePort: os.Getenv("HOSTNAME_PORT"),

		ApiKey:     os.Getenv("OPENAI_API_KEY"),
		SampleRate: 16_000,
		BitDepth:   16,
		Channels:   1,
		// jwt
		PrivateKey: os.Getenv("JWT_PRIVATE_KEY"),
		PublicKey:  os.Getenv("JWT_PUBLIC_KEY"),
		// db
		DBPath: os.Getenv("DB_PATH"),
	}
}
