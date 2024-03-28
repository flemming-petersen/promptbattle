package server

import (
	"github.com/flemming-petersen/promptbattle/config"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) showAdmin() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("admin", fiber.Map{})
	}
}

func (server *Server) openingRound() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState = OpeningState

		server.CurrentChallenge = &config.Challenge{
			Type:      "text",
			Challenge: "How would Flensburg look like in 100 years?",
		}
		server.PlayerPromptImages = map[string][]string{}
		server.PlayerFavorite = map[string]int{}

		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) announcingTheme() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState = AnnouncingState
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) startPrompting() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState = PromptState
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin?state=prompting")
	}
}

func (server *Server) generateImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState = GenerateState
		server.broadcastToAll(server.generateStateMsg())
		server.PlayerPromptImages = map[string][]string{}

		server.PlayerPromptImages["1"] = []string{
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
		}
		server.PlayerPromptImages["2"] = []string{
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
		}

		return c.Redirect("/admin")
	}
}

func (server *Server) startPickImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState = PickState
		server.PlayerFavorite = map[string]int{}
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}

func (server *Server) showFinalImages() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server.GameState = FinalState
		server.broadcastToAll(server.generateStateMsg())
		return c.Redirect("/admin")
	}
}
