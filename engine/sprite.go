package engine

import(
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type MultiSprite struct {
	sprites []*ebiten.Image

	rect image.Rectangle

	X float64
	Y float64
	Velx float64
	Vely float64

	Flip bool
}

func (this MultiSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if this.Flip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(32, 0)
	}

	op.GeoM.Translate(this.X, this.Y)

	// flip it
	

	for _, sprite := range this.sprites {
		screen.DrawImage(sprite.SubImage(this.rect).(*ebiten.Image), op)
	}
}
