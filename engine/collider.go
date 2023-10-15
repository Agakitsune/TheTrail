package engine

import (
	"os"
	"strconv"
	"strings"

	"fmt"

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
	Ground CollisionType = 9
	Death = 32
	Platform = 33
	Trigger = 34
)

type CollisionBox struct {
	Rect Rectangle
	Type CollisionType
	// Callback func()
}

type Collider struct {
	Boxes []CollisionBox

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
	type_ CollisionType
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

func (this *Collider) Update(game *Game, dood* MultiSprite) {
	for _, b := range this.Boxes {
		if (dood.Vely != 0) {
			println("velx: ", fmt.Sprintf("%f", dood.Velx))
			println("gx: ", fmt.Sprintf("%f", dood.X))
			println("bx: ", fmt.Sprintf("%f", b.Rect.X))
			println("bx: ", fmt.Sprintf("%f", b.Rect.X + b.Rect.Width))
			println("estimateMin: ", fmt.Sprintf("%f", (dood.X + 32 + dood.Velx)))
			println("estimateMax: ", fmt.Sprintf("%f", (dood.X - 32 + dood.Velx)))
			if ((dood.X + 22 + dood.Velx) <= float64(b.Rect.X + this.X) || (dood.X + 11 + dood.Velx) >= float64(b.Rect.X + b.Rect.Width + this.X)) {
				continue
			}
			println("vely: ", fmt.Sprintf("%f", dood.Vely))
			println("gy: ", fmt.Sprintf("%f", dood.Y))
			println("by: ", fmt.Sprintf("%f", b.Rect.Y))
			println("by + bh: ", fmt.Sprintf("%f",b.Rect.Y + b.Rect.Height))
			println("estimateMin: ", fmt.Sprintf("%f", (dood.Y + 32 + dood.Vely)))
			println("estimateMax: ", fmt.Sprintf("%f", (dood.Y - 32 + dood.Vely)))
			if ((dood.Y + 32 + dood.Vely) >= float64(b.Rect.Y + this.Y) && (dood.Y + 7 + dood.Vely) <= float64(b.Rect.Y + b.Rect.Height + this.Y)) {
				if b.Type == Death {
					println("YOU ARE DED, NOT BIG SURPRISE")
					continue
				} else if b.Type == Trigger {
					println("TRIGGER")
					continue
				}
				if dood.Vely > 0 {
					dood.Jump = false
					// dood.Y = float64(b.Rect.Y + this.Y - 32)
					dood.Airborne = false
				} else {
					// dood.Y = float64(b.Rect.Y + b.Rect.Height + this.Y)
				}
				dood.Vely = 0
			}
		}
		if (dood.Velx != 0) {
			// println("vely: ", fmt.Sprintf("%f", dood.Vely))
			// println("gy: ", fmt.Sprintf("%f", dood.Y))
			// println("by: ", fmt.Sprintf("%f", b.Rect.Y))
			// println("by + bh: ", fmt.Sprintf("%f",b.Rect.Y + b.Rect.Height))
			// println("estimateMin: ", fmt.Sprintf("%f", (dood.Y + 32 + dood.Vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (dood.Y - 32 + dood.Vely)))
			if ((dood.Y + 32 + dood.Vely) <= float64(b.Rect.Y + this.Y) || (dood.Y + 7 + dood.Vely) >= float64(b.Rect.Y + b.Rect.Height + this.Y)) {
				continue
			}
			// println("velx: ", fmt.Sprintf("%f", dood.Velx))
			// println("gx: ", fmt.Sprintf("%f", dood.X))
			// println("bx: ", fmt.Sprintf("%f", b.Rect.X))
			// println("bx: ", fmt.Sprintf("%f", b.Rect.X + b.Rect.Width))
			// println("estimateMin: ", fmt.Sprintf("%f", (dood.X + 32 + dood.Velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (dood.X - 32 + dood.Velx)))
			if ((dood.X + 22 + dood.Velx) >= float64(b.Rect.X + this.X) && (dood.X + 11 + dood.Velx) <= float64(b.Rect.X + b.Rect.Width + this.X)) {
				if b.Type == Death {
					println("YOU ARE DED, NOT BIG SURPRISE")
					continue
				} else if b.Type == Trigger {
					println("TRIGGER")
					continue
				}
				// else if b.Type == Scene {
				// 	println("Scene")
				// 	continue
				// }
				if dood.Velx > 0 {
					// dood.X = float64(b.Rect.X + this.X - 22)
				} else {
					// dood.X = float64(b.Rect.X + b.Rect.Width + this.X - 11)
				}
				dood.Velx = 0
			}
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
}
