package engine

import (
	"github.com/hajimehoshi/ebiten/v2"

	camera "github.com/melonfunction/ebiten-camera"
)

type State interface {
	Load(*Game)
	Update() error
	Draw(*ebiten.Image, *camera.Camera)
}
