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
		server.GameState.SetPhaseOpening()

		server.CurrentChallenge = &config.Challenge{
			Type:      "text",
			Challenge: "How would Flensburg look like in 100 years?",
		}

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

		server.GameState.SetImages("1", []string{
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
		})
		server.GameState.SetImages("2", []string{
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
			"https://placehold.co/600x400",
		})

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
