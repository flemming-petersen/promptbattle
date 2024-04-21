package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) startWebsocket(conn *websocket.Conn, msgHandler func(server *Server, msg map[string]interface{})) {
	// create outbound channel for this websocket connection
	server.FrontendMessages[conn] = make(chan []byte)
	// close channel
	defer close(server.FrontendMessages[conn])
	// delete closed channel from map with open frontendMessages websockets
	defer delete(server.FrontendMessages, conn)

	// send asny messages
	go func(frontendMessages chan []byte) {
		// read own frontendMessages channel and send to own websocket
		for msg := range frontendMessages {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}(server.FrontendMessages[conn])

	// read incoming messages, send to all other connection frontendMessages channels
	var (
		websocketMessageType int
		rawMsg               []byte
		err                  error
	)
	for {
		// read message
		websocketMessageType, rawMsg, err = conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("read:", err)
			}

			break
		}

		// handle only TextMessages
		if websocketMessageType != websocket.TextMessage {
			continue
		}

		// parse json
		msg := map[string]interface{}{}
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			fmt.Println(err)
			break
		}

		msgHandler(server, msg)
	}
}

func (server *Server) beamerWebsocket() func(*fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		server.startWebsocket(conn, func(server *Server, msg map[string]interface{}) {
			messageType := msg["type"].(string)
			switch messageType {
			case "ready":
				// player connection is ready!
				server.send(conn, server.generateStateMsg())
			default:
				fmt.Printf("Unknown message type %s from beamer\n", messageType)
			}
		})
	})
}

// playerWebsocket is the player controller
func (server *Server) playerWebsocket() func(*fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		server.startWebsocket(conn, func(server *Server, msg map[string]interface{}) {
			playerID := conn.Params("id")

			if !server.GameState.IsPlayerExist(playerID) {
				fmt.Printf("[Player: %s] Player does not exist\n", playerID)
				return
			}

			messageType := msg["type"].(string)
			switch messageType {
			case "ready":
				fmt.Printf("[Player: %s] Ready connected!\n", playerID)

				// player connection is ready!
				server.send(conn, server.generateStateMsg())
			case "prompt":
				// input message from player, add player before broadcast
				msg["player"] = playerID

				// set prompt for player
				server.GameState.SetPrompt(playerID, msg["prompt"].(string))

				// broadcast to all players
				server.broadcastToOther(conn, msg)
			case "pick":
				image := int(msg["image"].(float64))

				fmt.Printf("[Player: %s] select image %d\n", playerID, image)

				// set favorite image for player
				server.GameState.SetFavoriteImage(playerID, image)

				if server.GameState.SendAllPlayersFavoriteImage() {
					server.GameState.SetPhaseFinal()
					server.broadcastToAll(server.generateStateMsg())
				}
			default:
				fmt.Printf("[Player: %s] Unknown message type %s\n", playerID, messageType)
			}
		})
	})
}
