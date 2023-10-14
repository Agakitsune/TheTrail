package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type State interface {
	Load(*Game)
	Update() error
	Draw(*ebiten.Image)
}
