package states

import (
	"TheTrail/engine"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

type WorldMapState struct {
	game *engine.Game
}

func (s *WorldMapState) Load(g *engine.Game) {
	s.game = g
	fmt.Println("WorldMapState Load")
}

func (s *WorldMapState) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		fmt.Println("WorldMapState Update")
		s.game.SetState(&PlayState{})
	}
}

func (s *WorldMapState) Draw() {
	// fmt.Println("WorldMapState Draw")
}
