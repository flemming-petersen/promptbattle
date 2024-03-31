package server

type GameState string

const (
	// moderator introduce, player walk's to workstations
	OpeningState = GameState("opening")
	// theme reveal
	AnnouncingState = GameState("announcing")
	// players write prompts
	PromptState = GameState("prompt")
	// openai work
	GenerateState = GameState("generate")
	// select your best
	PickState = GameState("pick")
	// show final picked images in large
	FinalState = GameState("final")
)

func (server *Server) generateStateMsg() map[string]interface{} {
	stateMsg := map[string]interface{}{
		"type":          "state",
		"state":         server.GameState,
		"playerPrompts": server.PlayerPrompts,
	}
	if server.GameState != OpeningState {
		stateMsg["challenge"] = server.CurrentChallenge
	}
	if server.GameState == PickState || server.GameState == FinalState {
		stateMsg["playerImages"] = server.PlayerPromptImages
	}
	if server.GameState == FinalState {
		stateMsg["playerFavorite"] = server.PlayerFavorite
	}

	return stateMsg
}
