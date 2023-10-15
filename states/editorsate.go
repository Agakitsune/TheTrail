package states

import (
	"TheTrail/engine"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

type EditorState struct {
	game *engine.Game
}

func (s *EditorState) Load(g *engine.Game) {
	s.game = g
	fmt.Println("EditorState Load")
}

func (s *EditorState) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		fmt.Println("EditorState Update")
		s.game.SetState(&PlayState{})
	}
	return nil;
}

func (s *EditorState) Draw(screen *ebiten.Image) {
	// fmt.Println("EditorState Draw")
}
