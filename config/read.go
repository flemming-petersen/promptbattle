package config

import (
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
)

func ReadConfig() *Config {
	fileContent, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	config := &Config{}

	if err := json.Unmarshal(fileContent, &config.Challenges); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
        panic(err)
    }

	config.OpenAiKey = os.Getenv("OPENAI_KEY")

	return config
}
