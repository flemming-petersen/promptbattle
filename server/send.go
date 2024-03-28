package server

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func (server *Server) broadcastToAll(msg map[string]interface{}) {
	rawMsg, _ := json.Marshal(msg)

	// websocket connection, outbound message channel
	for _, frontendMessages := range server.FrontendMessages {
		// add new message to frontend channel from other connection
		frontendMessages <- rawMsg
	}
}

func (server *Server) broadcastToOther(conn *websocket.Conn, msg map[string]interface{}) {
	rawMsg, _ := json.Marshal(msg)

	// websocket connection, outbound message channel
	for otherConn, frontendMessages := range server.FrontendMessages {
		// ignore own connection
		if otherConn == conn {
			continue
		}

		// add new message to frontend channel from other connection
		frontendMessages <- rawMsg
	}
}

func (server *Server) send(conn *websocket.Conn, msg map[string]interface{}) {
	rawMsg, _ := json.Marshal(msg)

	if channel, ok := server.FrontendMessages[conn]; ok {
		channel <- rawMsg
	}
}
