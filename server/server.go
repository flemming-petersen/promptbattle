package server

import (
	"fmt"

	"github.com/flemming-petersen/promptbattle/config"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
)

type Server struct {
	App *fiber.App

	GameState          GameState
	CurrentChallenge   *config.Challenge
	PlayerPromptImages map[string][]string
	PlayerFavorite     map[string]int

	FrontendMessages map[*websocket.Conn]chan []byte
}

func NewServer() *Server {
	engine := html.New("./views", ".html")
	engine.Reload(true)

	server := &Server{
		App: fiber.New(fiber.Config{
			Views: engine,
		}),
		GameState:          OpeningState,
		PlayerPromptImages: map[string][]string{},
		PlayerFavorite:     map[string]int{},
		FrontendMessages:   make(map[*websocket.Conn]chan []byte),
	}

	server.App.Get("/player/:id", func(c *fiber.Ctx) error {
		err := c.Render("player", fiber.Map{
			"ID": c.Params("id"),
		})
		if err != nil {
			fmt.Println(err)
		}
		return err
	})
	server.App.Get("/player/:id/ws", server.playerWebsocket())

	server.App.Get("/beamer", func(c *fiber.Ctx) error {
		return c.Render("beamer", fiber.Map{})
	})

	server.App.Get("/beamer/ws", server.beamerWebsocket())

	server.App.Get("/admin", server.showAdmin())
	server.App.Get("/admin/state/opening", server.openingRound())
	server.App.Get("/admin/state/announcing", server.announcingTheme())
	server.App.Get("/admin/state/prompt", server.startPrompting())
	server.App.Get("/admin/state/generate", server.generateImages())
	server.App.Get("/admin/state/pick", server.startPickImages())
	server.App.Get("/admin/state/final", server.showFinalImages())

	return server
}

func (server *Server) Run() error {
	return server.App.Listen(":3000")
}
