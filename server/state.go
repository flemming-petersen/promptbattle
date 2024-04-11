package server

import "github.com/flemming-petersen/promptbattle/models"

func (server *Server) generateStateMsg() map[string]interface{} {
	players := map[string]interface{}{}
	for _, player := range server.GameState.Players() {
		players[player.ID] = map[string]interface{}{
			"prompt":          player.Prompt,
			"generatedImages": player.GeneratedImages,
			"favoriteImage":   player.FavoriteImage,
			"color": 		   player.Color,
		}
	}

	stateMsg := map[string]interface{}{
		"type":    "state",
		"phase":   server.GameState.Phase(),
		"players": players,
	}

	if server.GameState.Phase() != models.OpeningState {
		stateMsg["challenge"] = server.CurrentChallenge
	}

	return stateMsg
}
