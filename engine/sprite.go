package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type MultiSprite struct {
	Sprites []*ebiten.Image

	Rect image.Rectangle

	X     float64
	Y     float64
	Velx  float64
	Vely  float64
	Scale Vector2

	Flip bool

	Initialized bool
}

func (this MultiSprite) Draw(screen *ebiten.Image) {
	if !this.Initialized {
		this.Scale = Vector2{X: 1, Y: 1}
		this.Initialized = true
	}

	op := &ebiten.DrawImageOptions{}

	if this.Flip {
		op.GeoM.Scale(-this.Scale.X, this.Scale.Y)
		op.GeoM.Translate(32, 0)
	} else {
		op.GeoM.Scale(this.Scale.X, this.Scale.Y)
	}

	op.GeoM.Translate(this.X, this.Y)

	// flip it

	for _, sprite := range this.Sprites {
		screen.DrawImage(sprite.SubImage(this.Rect).(*ebiten.Image), op)
	}
}
