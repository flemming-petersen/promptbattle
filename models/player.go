package models

type Player struct {
	ID              string
	Color           string
	Prompt          string
	GeneratedImages []string
	FavoriteImage   *int
}

func (p *Player) Clear() {
	p.Prompt = ""
	p.GeneratedImages = nil
	p.FavoriteImage = nil
}
