package engine

import (
	"os"
	"strconv"
	"strings"

	// "fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"image/color"
)

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

type CollisionType int

const (
	Ground   CollisionType = 9
	Death                  = 32
	Platform               = 33
	Trigger                = 34
)

type CollisionBox struct {
	Rect Rectangle
	Type CollisionType
	// Callback func()
}

type Collider struct {
	Boxes       []CollisionBox
	playerBoxes []Rectangle

	X int
	Y int
}

type lineInfo struct {
	x     int
	y     int
	width int
	type_ CollisionType
}

type rectInfo struct {
	x      int
	y      int
	width  int
	height int
	type_  CollisionType
}

func (r Rectangle) Collides(r2 Rectangle) bool {
	return r.X < r2.X+r2.Width &&
		r.X+r.Width > r2.X &&
		r.Y < r2.Y+r2.Height &&
		r.Y+r.Height > r2.Y
}

func (r Rectangle) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(r.X),
		float32(r.Y),
		float32(r.Width),
		float32(r.Height),
		color.RGBA{255, 0, 0, 255},
		false)
}

func NewColliderMap(path string, x, y int) *Collider {
	var collider *Collider = new(Collider)
	data, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	found := false
	pos := 0
	width := 0
	type_ := 0

	info := make([]lineInfo, 0)

	lines := strings.Split(strings.Replace(string(data), "\r", "", -1), "\n")
	for i, line := range lines {
		if len(line) > 1 {
			for j, el := range strings.Split(line, ",") {
				val, err := strconv.Atoi(el)

				if err != nil {
					panic(err)
				}

				if val != -1 {
					if found && val != type_ {
						info = append(info, lineInfo{pos, i, width, CollisionType(type_)})
						pos = j
						found = true
						type_ = val
						width = 1
					} else {
						if !found {
							pos = j
							found = true
							type_ = val
						}
						width++
					}
				} else {
					if found {
						info = append(info, lineInfo{pos, i, width, CollisionType(type_)})
						found = false
						width = 0
					}
				}
			}
			if found {
				info = append(info, lineInfo{pos, i, width, CollisionType(type_)})
			}
			found = false
			width = 0
		}
	}

	pending := make([]rectInfo, 0)
	for _, line := range info {
		newRect := true
		// println(line.x, line.y, line.width, line.type_)
		for i, rect := range pending {
			if rect.x == line.x && line.y == rect.y+rect.height {
				pending[i].height++
				newRect = false
				break
			}
		}

		if newRect {
			pending = append(pending, rectInfo{line.x, line.y, line.width, 1, line.type_})
		}
	}

	for _, rect := range pending {
		collider.Boxes = append(collider.Boxes,
			CollisionBox{
				Rectangle{rect.x * 8, rect.y * 8, rect.width * 8, rect.height * 8},
				rect.type_,
			},
		)
	}

	collider.X = x
	collider.Y = y

	return collider
}

func (this *Collider) Update(game *Game, dood *MultiSprite) {

	w, h := int(12), int(16)
	_ = w
	_ = h

	lp := []int{int(dood.X - dood.Velx), int(dood.Y - dood.Vely)}

	this.playerBoxes = make([]Rectangle, 4)
	this.playerBoxes[0] = Rectangle{int(dood.X) + 9, int(dood.Y) + (32-h)/2 + 4, 4, h}  // Left
	this.playerBoxes[1] = Rectangle{int(dood.X) + 19, int(dood.Y) + (32-h)/2 + 4, 4, h} // Right
	this.playerBoxes[2] = Rectangle{int(dood.X) + 10, int(dood.Y + 6), w, 4}            // Top
	this.playerBoxes[3] = Rectangle{int(dood.X) + 10, int(dood.Y + 30), w, 4}           // Bottom

	for _, b := range this.Boxes {
		// println("\nX: ", fmt.Sprintf("%f", dood.X))
		// println("Y: ", fmt.Sprintf("%f\n", dood.Y))

		rect := Rectangle{b.Rect.X + this.X, b.Rect.Y + this.Y, b.Rect.Width, b.Rect.Height}

		if this.playerBoxes[0].Collides(rect) && dood.Velx < 0 {
			dood.Velx = 0
			dood.X = float64(lp[0] - 1)
		}

		if this.playerBoxes[1].Collides(rect) && dood.Velx > 0 {
			dood.Velx = 0
			dood.X = float64(lp[0] + 1)
		}

		if this.playerBoxes[2].Collides(rect) && dood.Vely < 0 {
			dood.Vely = 0
			dood.Y = float64(lp[1] + 1)
		}

		if this.playerBoxes[3].Collides(rect) && dood.Vely > 0 {
			dood.Vely = 0
			dood.Y = float64(rect.Y - 32)
			dood.Airborne = false
			dood.Jump = false
		}
	}
}

func (this *Collider) Draw(screen *ebiten.Image) {
	for _, b := range this.Boxes {
		clr := color.RGBA{0, 0, 0, 255}
		if b.Type == Ground {
			clr = color.RGBA{0, 255, 255, 2}
		} else if b.Type == Death {
			clr = color.RGBA{255, 0, 0, 2}
		} else if b.Type == Platform {
			clr = color.RGBA{0, 255, 0, 2}
		} else if b.Type == Trigger {
			clr = color.RGBA{255, 0, 255, 2}
		} else if b.Type == Trigger {
			clr = color.RGBA{255, 255, 0, 2}
		}

		vector.DrawFilledRect(
			screen,
			float32(b.Rect.X),
			float32(b.Rect.Y),
			float32(b.Rect.Width),
			float32(b.Rect.Height),
			clr,
			false,
		)
	}

	for _, b := range this.playerBoxes {
		b.Draw(screen)
	}
}
