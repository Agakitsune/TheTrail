package engine

import(
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	camera "github.com/melonfunction/ebiten-camera"
)

type MultiSprite struct {
	sprites []*ebiten.Image

	rect image.Rectangle

	Jump bool
	Airborne bool
		
	X float64
	Y float64
	Velx float64
	Vely float64

	Flip bool
	Walling bool
}

func NewSprite(
	head, torso, boot string, rect image.Rectangle,
)* MultiSprite {
	var sprite* MultiSprite = new(MultiSprite)

	sprite.sprites = []*ebiten.Image{
		LoadImage(head),
		LoadImage(torso),
		LoadImage(boot),
	}

	sprite.rect = rect

	return sprite
}

func (this MultiSprite) Draw(screen *ebiten.Image, camera* camera.Camera) {
	op := &ebiten.DrawImageOptions{}

	if this.Flip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(32, 0)
	}

	op.GeoM.Translate(this.X, this.Y)
	op = camera.GetTranslation(op, 0, 0)

	// flip it
	

	for _, sprite := range this.sprites {
		screen.DrawImage(sprite.SubImage(this.rect).(*ebiten.Image), op)
	}
}
