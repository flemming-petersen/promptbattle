package models

type GamePhase string

const (
	// moderator introduce, player walk's to workstations
	OpeningState = GamePhase("opening")
	// theme reveal
	AnnouncingState = GamePhase("announcing")
	// players write prompts
	PromptingState = GamePhase("prompting")
	// openai work
	GenerateState = GamePhase("generate")
	// select your best
	PickingState = GamePhase("picking")
	// show final picked images in large
	FinalState = GamePhase("final")
)

type State struct {
	phase   GamePhase
	players map[string]*Player
}

func NewState() *State {
	return &State{
		phase: OpeningState,
		players: map[string]*Player{
			// only here defined playerIDs are allowed to join
			"1": {ID: "1", Color: "#00D6C4"},
			"2": {ID: "2", Color: "#d98524"},
		},
	}
}

func (s *State) IsPlayerExist(id string) bool {
	_, ok := s.players[id]
	return ok
}

func (s *State) Phase() GamePhase {
	return s.phase
}

func (s *State) Players() map[string]*Player {
	return s.players
}

func (s *State) SetPhaseOpening() {
	s.phase = OpeningState
	for _, player := range s.players {
		player.Clear()
	}
}

func (s *State) SetPhaseAnnouncing() {
	s.phase = AnnouncingState
}

func (s *State) SetPhasePrompting() {
	s.phase = PromptingState

	for _, player := range s.players {
		player.Prompt = ""
	}
}

func (s *State) SetPhaseGenerate() {
	s.phase = GenerateState

	for _, player := range s.players {
		player.GeneratedImages = nil
		player.FavoriteImage = nil
	}
}

func (s *State) SetPhasePicking() {
	s.phase = PickingState

	for _, player := range s.players {
		player.FavoriteImage = nil
	}
}

func (s *State) SetPhaseFinal() {
	s.phase = FinalState
}

func (s *State) SetImages(playerID string, images []string) {
	s.players[playerID].GeneratedImages = images
}

func (s *State) SetPrompt(playerID, prompt string) {
	s.players[playerID].Prompt = prompt
}

func (s *State) SetFavoriteImage(playerID string, image int) {
	s.players[playerID].FavoriteImage = &image
}

func (s *State) SendAllPlayersFavoriteImage() bool {
	for _, player := range s.players {
		if player.FavoriteImage == nil {
			return false
		}
	}

	return true
}
