package states

import (
	"TheTrail/engine"
	"fmt"
)

type PlayState struct {
	game *engine.Game
}

func (s *PlayState) Load(g *engine.Game) {
	s.game = g
	fmt.Println("PlayState Load")
}

func (s *PlayState) Update() {
	fmt.Println("PlayState Update")
}

func (s *PlayState) Draw() {
	fmt.Println("PlayState Draw")
}
