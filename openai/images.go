package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/flemming-petersen/promptbattle/config"
)

type Client struct {
	Config *config.Config
	APIURL string
}

func NewClient(configuration *config.Config) *Client {
	return &Client{
		Config: configuration,
		APIURL: "https://api.openai.com/v1/images/generations",
	}
}

type requestImagePayload struct {
	Model string `json:"model"`
	Prompt string `json:"prompt"`
	N int `json:"n"`
	Size string `json:"size"`
}

type responseImagePayload struct {
	Images []struct{
		URL string `json:"url"`
	} `json:"data"`
}

func (client *Client) GeneratedImages(prompt string) ([]string, error) {
	payload := requestImagePayload{
		Model: "dall-e-2",
		Prompt: prompt,
		N: 4,
		Size: "256x256",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", client.APIURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+client.Config.OpenAiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// error!
		responsePayload, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		fmt.Println("error response", string(responsePayload))
		return nil, errors.New("error from openAI")
	}

	responsePayload := &responseImagePayload{}
	err = json.NewDecoder(response.Body).Decode(responsePayload)
	if err != nil {
		return nil, err
	}

	images := []string{}
	for _, image := range responsePayload.Images {
		images = append(images, image.URL)
	}

	return images, nil
}
