package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (server *Server) showAdmin() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("admin", fiber.Map{})
	}
}

func (server *Server) openingRound() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState.SetPhaseOpening()

		server.CurrentChallengeIndex++

		if server.CurrentChallengeIndex >= len(server.Config.Challenges) {
			server.CurrentChallengeIndex = 0
		}

		server.CurrentChallenge = server.Config.Challenges[server.CurrentChallengeIndex]

		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) announcingTheme() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState.SetPhaseAnnouncing()

		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) startPrompting() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState.SetPhasePrompting()

		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) generateImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState.SetPhaseGenerate()
		server.broadcastToAll(server.generateStateMsg())


		// Send prompt to openai
		for playerID, player := range server.GameState.Players() {
			fmt.Printf("[Player: %s] Prompt: %s\n", playerID, player.Prompt)

			images, err := server.OpenAiClient.GeneratedImages(player.Prompt)
			if err != nil {
				return err
			}

			server.GameState.SetImages(playerID, images)
		}

		server.GameState.SetPhasePicking()
		server.broadcastToAll(server.generateStateMsg())

		return c.Redirect("/admin")
	}
}

func (server *Server) startPickImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState.SetPhasePicking()
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) showFinalImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState.SetPhaseFinal()
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}
