package engine

import(
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type MultiSprite struct {
	sprites []*ebiten.Image

	rect image.Rectangle

	x float64
	y float64
	velx float64
	vely float64

	flip bool
}

func (this MultiSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if this.flip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(32, 0)
	}

	op.GeoM.Translate(this.x, this.y)

	// flip it
	

	for _, sprite := range this.sprites {
		screen.DrawImage(sprite.SubImage(this.rect).(*ebiten.Image), op)
	}
}
