package engine

import (
	"github.com/hajimehoshi/ebiten/v2/audio"

	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
)

type Music struct {
	player *audio.Player
	context *audio.Context
}


