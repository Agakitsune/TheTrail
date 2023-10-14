package engine

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"io/ioutil"

	// "fmt"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 180

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

var (
	runnerImage *ebiten.Image
	tileset     *ebiten.Image
)

func LoadGaYme() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)

	// Loading tileset
	content, error := ioutil.ReadFile("./assets/tileset.png")
	if error != nil {
		log.Fatal(error)
	}

	img, _, err = image.Decode(bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	tileset = ebiten.NewImageFromImage(img)
}

type Rectangle struct {
	x float32
	y float32
	w float32
	h float32
}

type Game struct {
	count int

	inited bool

	dood MultiSprite

	keys []ebiten.Key

	rects []Rectangle

	jump bool
	airborne bool

	tilemap *Tilemap
	animator *Animator

	state State
}

func (r Rectangle) Draw(screen *ebiten.Image, color color.RGBA) {
	vector.DrawFilledRect(screen, r.x, r.y, r.w, r.h, color, false)
}

func (g *Game) SetState(state State) {
	g.state = state
	g.state.Load(g)
}

func (g *Game) Init() {
	g.rects = []Rectangle{
		Rectangle{0, 20 * 8, 40 * 8, 3 * 8},
		Rectangle{33 * 8, 16 * 8, 7 * 8, 4 * 8},
		Rectangle{23 * 8, 16 * 8, 2 * 8, 4 * 8},
		Rectangle{28 * 8, 14 * 8, 2 * 8, 2 * 8},
	}
	g.inited = true
	g.jump = false;

	g.dood = MultiSprite{
		sprites: []*ebiten.Image{
			LoadImage("./assets/dood/boots_one.png"),
			LoadImage("./assets/dood/torso_three.png"),
			LoadImage("./assets/dood/head_two.png"),
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

	g.tilemap = NewTilemap("./assets/map.csv", "./assets/tileset.png")

	// Initialize the state
	g.state = nil
}

func (g *Game) Update() error {
	if !g.inited {
		g.Init()
	}
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count++

	g.animator.Update(&g.dood)
	if g.state != nil {
		g.state.Update()
	}
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
			if ((g.dood.x + 22 + g.dood.velx) <= float64(b.x) || (g.dood.x + 11 + g.dood.velx) >= float64(b.x + b.w)) {
				continue
			}
			// println("vely: ", fmt.Sprintf("%f", g.dood.vely))
			// println("gy: ", fmt.Sprintf("%f", g.dood.y))
			// println("by: ", fmt.Sprintf("%f", b.y))
			// println("by + bh: ", fmt.Sprintf("%f",b.y + b.h))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.dood.y + 32 + g.dood.vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.dood.y - 32 + g.dood.vely)))
			if ((g.dood.y + 32 + g.dood.vely) >= float64(b.y) && (g.dood.y + 7 + g.dood.vely) <= float64(b.y + b.h)) {
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
			if ((g.dood.y + 32 + g.dood.vely) <= float64(b.y) || (g.dood.y + 7 + g.dood.vely) >= float64(b.y + b.h)) {
				continue
			}
			// println("velx: ", fmt.Sprintf("%f", g.dood.velx))
			// println("gx: ", fmt.Sprintf("%f", g.dood.x))
			// println("bx: ", fmt.Sprintf("%f", b.x))
			// println("bx: ", fmt.Sprintf("%f", b.x + b.w))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.dood.x + 32 + g.dood.velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.dood.x - 32 + g.dood.velx)))
			if ((g.dood.x + 22 + g.dood.velx) >= float64(b.x) && (g.dood.x + 11 + g.dood.velx) <= float64(b.x + b.w)) {
				g.dood.velx = 0
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
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
	return ScreenWidth, ScreenHeight
}

func CreateGame(state State) *Game {
	var game = &Game{}
	game.Init()
	game.SetState(state)

	return game
}
