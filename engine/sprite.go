package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	camera "github.com/melonfunction/ebiten-camera"
)

type MultiSprite struct {
	Sprites []*ebiten.Image

	Rect image.Rectangle

	Jump bool
	Airborne bool
		
	X float64
	Y float64
	Velx float64
	Vely float64

	Flip bool

	TryClimb bool
	Climbing bool

	TrySlowFall bool
	SlowFall bool
	Dir int
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
		op.GeoM.Scale(-this.Scale.X, this.Scale.Y)
		op.GeoM.Translate(32, 0)
	} else {
		op.GeoM.Scale(this.Scale.X, this.Scale.Y)
	}

	op.GeoM.Translate(this.X, this.Y)
	op = camera.GetTranslation(op, 0, 0)

	// flip it

	for _, sprite := range this.Sprites {
		screen.DrawImage(sprite.SubImage(this.Rect).(*ebiten.Image), op)
	}
}
