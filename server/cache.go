package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// savePromptAndImages saves the prompt and images for a player in a round and returns the relative paths to the images.
func (server *Server) savePromptAndImages(generateID string, playerID string, prompt string, imageURLs []string) ([]string, error) {
	relativeBasePath := filepath.Join(generateID, playerID)

	// Create the folder for the player
	if err := os.MkdirAll(filepath.Join(server.Config.PromptImageBasePath, relativeBasePath), 0755); err != nil {
		return nil, err
	}

	// Save the prompt
	promptPath := filepath.Join(server.Config.PromptImageBasePath, relativeBasePath, "prompt.txt")
	if err := os.WriteFile(promptPath, []byte(prompt), 0644); err != nil {
		return nil, err
	}

	// Download and save the images asynchronously
	var waitGroup sync.WaitGroup
	for imageNumber, imageURL := range imageURLs {
		waitGroup.Add(1)

		go func(imageNumber int, imageURL string) {
			defer waitGroup.Done()

			imageResp, err := http.Get(imageURL)
			if err != nil {
				return
			}

			defer imageResp.Body.Close()

			cacheFilename := filepath.Join(server.Config.PromptImageBasePath, relativeBasePath, fmt.Sprintf("image_%d.jpg", imageNumber))
			file, err := os.Create(cacheFilename)
			if err != nil {
				return
			}

			defer file.Close()

			_, err = io.Copy(file, imageResp.Body)
			if err != nil {
				return
			}
		}(imageNumber, imageURL)
	}

	// Wait for all images to be downloaded
	waitGroup.Wait()

	// Return the paths to the images
	imagePaths := make([]string, len(imageURLs))
	for i := range imageURLs {
		imagePaths[i] = filepath.Join(server.Config.PromptImageURLPrefix, relativeBasePath, fmt.Sprintf("image_%d.jpg", i))
	}

	return imagePaths, nil
}
