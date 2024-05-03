package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (server *Server) showAdmin() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("admin", fiber.Map{})
	}
}

func (server *Server) openingRound() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("[Admin] Opening and reset round")

		server.GameState.SetPhaseOpening()

		server.CurrentChallengeIndex++

		if server.CurrentChallengeIndex >= len(server.Config.Challenges) {
			server.CurrentChallengeIndex = 0
		}

		fmt.Printf("[Admin] Set challenge to index %d\n", server.CurrentChallengeIndex)

		server.CurrentChallenge = server.Config.Challenges[server.CurrentChallengeIndex]

		fmt.Printf("[Admin] Set challenge to %s\n", server.CurrentChallenge.Challenge)

		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) announcingTheme() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("[Admin] Announce theme")

		server.GameState.SetPhaseAnnouncing()

		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) startPrompting() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("[Admin] Start prompting")

		server.GameState.SetPhasePrompting()

		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) generateImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("[Admin] Start generating images")

		server.GameState.SetPhaseGenerate()
		server.broadcastToAll(server.generateStateMsg())

		generateID := fmt.Sprintf("%d", time.Now().Unix())
		startTime := time.Now()

		// Send prompt to openai
		for playerID, player := range server.GameState.Players() {
			fmt.Printf("[Player: %s] Prompt: %s\n", playerID, player.Prompt)

			if player.Prompt == "" {
				fmt.Printf("[Player: %s] No prompt\n", playerID)
				player.Prompt = "matrix, many numbers, a room, emptiness, glitch effects, high quality, award winning"
			}

			// generate images
			imageURLs, err := server.OpenAiClient.GeneratedImages(player.Prompt)
			if err != nil {
				fmt.Printf("[Player: %s] Error: %s\n", playerID, err.Error())
				return err
			}

			imagePaths, err := server.savePromptAndImages(generateID, playerID, player.Prompt, imageURLs)
			if err != nil {
				return err
			}

			server.GameState.SetImages(playerID, imagePaths)
		}

		fmt.Printf("[Admin] Generate images took %s\n", time.Since(startTime))

		server.GameState.SetPhasePicking()
		server.broadcastToAll(server.generateStateMsg())

		return c.Redirect("/admin")
	}
}

func (server *Server) startPickImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("[Admin] Start picking images")

		server.GameState.SetPhasePicking()
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) showFinalImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("[Admin] Show final images")

		server.GameState.SetPhaseFinal()
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}
