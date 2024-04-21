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
	config.PromptImageBasePath = os.Getenv("IMAGE_BASE_PATH")
	config.PromptImageURLPrefix = os.Getenv("IMAGE_URL_PREFIX")

	if config.PromptImageBasePath == "" {
		config.PromptImageBasePath = "./images"
	}

	if config.PromptImageURLPrefix == "" {
		config.PromptImageURLPrefix = "/prompts"
	}

	return config
}
