// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"image"
	_ "image/png"
	"image/color"
	"log"

	"io/ioutil"

	// "fmt"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

const (
	screenWidth  = 320
	screenHeight = 180

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

var (
	runnerImage *ebiten.Image
	tileset *ebiten.Image
)

type Rectangle struct {
	x float32
	y float32
	w float32
	h float32
}

type MultiSprite struct {
	sprites []*ebiten.Image

	rect image.Rectangle

	x float64
	y float64
	velx float64
	vely float64

	flip bool
}

type Game struct {
	count int

	inited bool

	dood MultiSprite

	keys []ebiten.Key

	rects []Rectangle

	jump bool
	airborne bool

	tilemap Tilemap

	animator *Animator
}

type Tilemap struct {
	tileset *ebiten.Image

	tiles []int
}

type Animation struct {
	frames []int
	row int

	loopOn int

	selection int
	speed int

	count int
}

type Animator struct {
	animations map[string]*Animation

	current string
}

func (this *Animation) GetFrame(frame int) int {
	return this.frames[frame]
}

func (this *Animation) Update(sprite *MultiSprite) {
	this.count++

	if this.count % this.speed != 0 {
		return
	}
	this.selection++

	if this.selection >= len(this.frames) {
		this.selection = this.loopOn
	}

	frame := this.GetFrame(this.selection)
	sprite.rect = image.Rect(frame * 32, this.row * 32, frame * 32 + 32, this.row * 32 + 32)
}

func (this *Animator) SetAnimation(name string) {
	if this.current == name {
		return
	}
	var anim = this.animations[this.current]
	anim.selection = 0
	this.current = name
	anim = this.animations[this.current]
	anim.selection = 0
}

func (this *Animator) Update(sprite *MultiSprite) {
	var anim = this.animations[this.current]
	anim.Update(sprite)
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

func (r Rectangle) Draw(screen *ebiten.Image, color color.RGBA) {
	vector.DrawFilledRect(screen, r.x, r.y, r.w, r.h, color, false)
}

func (g *Game) init() {
	g.rects = []Rectangle{
		Rectangle{0, 20 * 8, 40 * 8, 3 * 8},
		Rectangle{33 * 8, 16 * 8, 7 * 8, 4 * 8},
	}
	g.inited = true
	g.jump = false;

	g.dood = MultiSprite{
		sprites: []*ebiten.Image{
			loadImage("./assets/dood/boots_one.png"),
			loadImage("./assets/dood/torso_three.png"),
			loadImage("./assets/dood/head_two.png"),
		},
		rect: image.Rect(0, 0, 32, 32),
		x: 0,
		y: 0,
		velx: 0,
		vely: 0,
	}

	g.animator = &Animator{
		animations: map[string]*Animation{
			"idle": &Animation{
				frames: []int{0, 1, 2, 3, 4, 5},
				row: 0,
				loopOn: 0,
				selection: 0,
				speed: 10,
			},
			"walk": &Animation{
				frames: []int{0, 1, 2, 3, 4, 5, 6},
				row: 1,
				loopOn: 1,
				selection: 0,
				speed: 10,
			},
			"run": &Animation{
				frames: []int{0, 1, 2, 3, 4, 5},
				row: 2,
				loopOn: 0,
				selection: 0,
				speed: 10,
			},
			"jump": &Animation{
				frames: []int{0, 1, 2},
				row: 3,
				loopOn: 2,
				selection: 0,
				speed: 5,
			},
		},
		current: "idle",
	}

	g.tilemap = Tilemap{
		tileset: tileset,
		tiles: []int{
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,0,1,1,1,1,1,1,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,6,3,4,4,4,4,4,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,6,9,7,7,7,7,7,
			-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,6,9,7,7,7,7,7,
			0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,18,9,7,7,7,7,7,
			6,3,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,24,7,7,7,7,7,
			6,9,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,
		},
	}
}

func (tilemap Tilemap) Draw(screen *ebiten.Image) {
	for i, tile := range tilemap.tiles {
		if tile != -1 {
			op := &ebiten.DrawImageOptions{}
			sx, sy := (tile % 6 * 8), tile / 6 * 8
			op.GeoM.Translate(float64(i % 40) * 8, float64(i / 40) * 8)
			screen.DrawImage(tilemap.tileset.SubImage(image.Rect(sx, sy, sx+8, sy+8)).(*ebiten.Image), op)
		}
	}
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	g.animator.Update(&g.dood)
	return nil
}

func (g *Game) collide(boxes []Rectangle) {
	for _, b := range boxes {
		if (g.dood.vely != 0) {
			// println("velx: ", fmt.Sprintf("%f", g.dood.velx))
			// println("gx: ", fmt.Sprintf("%f", g.dood.x))
			// println("bx: ", fmt.Sprintf("%f", b.x))
			// println("bx: ", fmt.Sprintf("%f", b.x + b.w))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.dood.x + 32 + g.dood.velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.dood.x - 32 + g.dood.velx)))
			if ((g.dood.x + 32 + g.dood.velx) <= float64(b.x) || (g.dood.x - 32 + g.dood.velx) >= float64(b.x + b.w)) {
				continue
			}
			// println("vely: ", fmt.Sprintf("%f", g.dood.vely))
			// println("gy: ", fmt.Sprintf("%f", g.dood.y))
			// println("by: ", fmt.Sprintf("%f", b.y))
			// println("by + bh: ", fmt.Sprintf("%f",b.y + b.h))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.dood.y + 32 + g.dood.vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.dood.y - 32 + g.dood.vely)))
			if ((g.dood.y + 32 + g.dood.vely) >= float64(b.y) && (g.dood.y - 32 + g.dood.vely) <= float64(b.y + b.h)) {
				if g.dood.vely > 0 {
					g.jump = false
					g.airborne = false
				}
				g.dood.vely = 0
			}
		}
		if (g.dood.velx != 0) {
			// println("vely: ", fmt.Sprintf("%f", g.dood.vely))
			// println("gy: ", fmt.Sprintf("%f", g.dood.y))
			// println("by: ", fmt.Sprintf("%f", b.y))
			// println("by + bh: ", fmt.Sprintf("%f",b.y + b.h))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.dood.y + 32 + g.dood.vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.dood.y - 32 + g.dood.vely)))
			if ((g.dood.y + 32 + g.dood.vely) <= float64(b.y) || (g.dood.y - 32 + g.dood.vely) >= float64(b.y + b.h)) {
				continue
			}
			// println("velx: ", fmt.Sprintf("%f", g.dood.velx))
			// println("gx: ", fmt.Sprintf("%f", g.dood.x))
			// println("bx: ", fmt.Sprintf("%f", b.x))
			// println("bx: ", fmt.Sprintf("%f", b.x + b.w))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.dood.x + 32 + g.dood.velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.dood.x - 32 + g.dood.velx)))
			if ((g.dood.x + 32 + g.dood.velx) >= float64(b.x) && (g.dood.x - 32 + g.dood.velx) <= float64(b.x + b.w)) {
				g.dood.velx = 0
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	// op.GeoM.Translate(g.dood.x, g.dood.y)
	// i := (g.count / 5) % frameCount
	// sx, sy := frameOX+i*frameWidth, frameOY
	// screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

	g.dood.Draw(screen)

	moveX := false
	run := true

	g.dood.vely += 0.1

	vel := 0.5
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		vel = 0.2
		run = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.dood.velx = vel
		g.dood.flip = false
		if !g.airborne {
			if run {
				g.animator.SetAnimation("run")
			} else {
				g.animator.SetAnimation("walk")
			}
		}
		moveX = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.dood.velx = -vel
		g.dood.flip = true
		if !g.airborne {
			if run {
				g.animator.SetAnimation("run")
			} else {
				g.animator.SetAnimation("walk")
			}
		}
		moveX = true
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.jump {
		g.animator.SetAnimation("jump")
		g.dood.vely = -3
		g.jump = true
		g.airborne = true
	}
	if !moveX {
		if !g.airborne {
			g.animator.SetAnimation("idle")
		}
		g.dood.velx = 0
	}

	g.collide(g.rects)

	g.dood.x += g.dood.velx
	g.dood.y += g.dood.vely

	g.tilemap.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func loadImage(path string) *ebiten.Image {
	content, error := ioutil.ReadFile(path)
	if error != nil {
		log.Fatal(error)
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

func main() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)

	tileset = loadImage("./assets/tileset.png")

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}