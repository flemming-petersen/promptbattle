package config

type Config struct {
	OpenAiKey  string
	Challenges []*Challenge

	PromptImageBasePath  string
	PromptImageURLPrefix string
}
