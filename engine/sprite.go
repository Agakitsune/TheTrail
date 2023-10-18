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
	Scale Vector2

	Flip bool

	TryClimb bool
	Climbing bool
	Edge bool

	Stamina int

	TrySlowFall bool
	SlowFall bool
	Dir int

	Falltrough bool

	Dead bool

	Win bool
}

func NewSprite(
	paths []string, rect image.Rectangle, scale Vector2,
)* MultiSprite {
	var sprite* MultiSprite = new(MultiSprite)

	sprite.Sprites = make([]*ebiten.Image, 0)

	for _, path := range paths {
		sprite.Sprites = append(sprite.Sprites, LoadImage(path))
	}

	sprite.Scale = scale

	sprite.Rect = rect

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
	if camera != nil {
		op = camera.GetTranslation(op, 0, 0)
	}

	// flip it

	for _, sprite := range this.Sprites {
		screen.DrawImage(sprite.SubImage(this.Rect).(*ebiten.Image), op)
	}
}
