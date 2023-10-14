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

	color color.RGBA
}

type Game struct {
	count int

	inited bool

	x    float64
	y    float64
	velx float64
	vely float64

	keys []ebiten.Key

	rects []Rectangle

	jump bool

	tilemap Tilemap

	state State
}

type Tilemap struct {
	tileset *ebiten.Image

	tiles []int
}

func (r Rectangle) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, r.x, r.y, r.w, r.h, r.color, false)
}

func (g *Game) SetState(state State) {
	g.state = state
	g.state.Load(g)
}

func (g *Game) Init() {
	g.inited = true
	// Rects
	g.rects = []Rectangle{
		Rectangle{0, 20 * 8, 40 * 8, 3 * 8, color.RGBA{255, 255, 255, 255}},
		Rectangle{33 * 8, 16 * 8, 7 * 8, 4 * 8, color.RGBA{255, 255, 255, 255}},
		// Rectangle{200, 220 - 50, 50, 50, color.RGBA{255, 255, 255, 255}},
		// Rectangle{100, 220 - 100, 50, 20, color.RGBA{255, 255, 255, 255}},
		// Rectangle{250, 0, 70, 240, color.RGBA{255, 255, 255, 255}},
	}
	g.inited = true
	g.jump = false

	g.tilemap = Tilemap{
		tileset: tileset,
		tiles: []int{
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0, 1, 1, 1, 1, 1, 1,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 6, 3, 4, 4, 4, 4, 4,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 6, 9, 7, 7, 7, 7, 7,
			-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 6, 9, 7, 7, 7, 7, 7,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 18, 9, 7, 7, 7, 7, 7,
			6, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 24, 7, 7, 7, 7, 7,
			6, 9, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		},
	}

	// Initialize the state
	g.state = nil
}

func (tilemap Tilemap) Draw(screen *ebiten.Image) {
	for i, tile := range tilemap.tiles {
		if tile != -1 {
			op := &ebiten.DrawImageOptions{}
			sx, sy := (tile % 6 * 8), tile/6*8
			op.GeoM.Translate(float64(i%40)*8, float64(i/40)*8)
			screen.DrawImage(tilemap.tileset.SubImage(image.Rect(sx, sy, sx+8, sy+8)).(*ebiten.Image), op)
		}
	}
}

func (g *Game) Update() error {
	if !g.inited {
		g.Init()
	}
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count++

	if g.state != nil {
		g.state.Update()
	}
	return nil
}

func (g *Game) collide(boxes []Rectangle) {
	for _, b := range boxes {
		if g.vely != 0 {
			// println("velx: ", fmt.Sprintf("%f", g.velx))
			// println("gx: ", fmt.Sprintf("%f", g.x))
			// println("bx: ", fmt.Sprintf("%f", b.x))
			// println("bx: ", fmt.Sprintf("%f", b.x + b.w))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.x + 16 + g.velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.x - 16 + g.velx)))
			if (g.x+16+g.velx) <= float64(b.x) || (g.x-16+g.velx) >= float64(b.x+b.w) {
				continue
			}
			// println("vely: ", fmt.Sprintf("%f", g.vely))
			// println("gy: ", fmt.Sprintf("%f", g.y))
			// println("by: ", fmt.Sprintf("%f", b.y))
			// println("by + bh: ", fmt.Sprintf("%f",b.y + b.h))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.y + 16 + g.vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.y - 16 + g.vely)))
			if (g.y+16+g.vely) >= float64(b.y) && (g.y-16+g.vely) <= float64(b.y+b.h) {
				if g.vely > 0 {
					g.jump = false
				}
				g.vely = 0
				b.color = color.RGBA{255, 255, 0, 255}
			}
		}
		if g.velx != 0 {
			// println("vely: ", fmt.Sprintf("%f", g.vely))
			// println("gy: ", fmt.Sprintf("%f", g.y))
			// println("by: ", fmt.Sprintf("%f", b.y))
			// println("by + bh: ", fmt.Sprintf("%f",b.y + b.h))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.y + 16 + g.vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.y - 16 + g.vely)))
			if (g.y+16+g.vely) <= float64(b.y) || (g.y-16+g.vely) >= float64(b.y+b.h) {
				continue
			}
			// println("velx: ", fmt.Sprintf("%f", g.velx))
			// println("gx: ", fmt.Sprintf("%f", g.x))
			// println("bx: ", fmt.Sprintf("%f", b.x))
			// println("bx: ", fmt.Sprintf("%f", b.x + b.w))
			// println("estimateMin: ", fmt.Sprintf("%f", (g.x + 16 + g.velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (g.x - 16 + g.velx)))
			if (g.x+16+g.velx) >= float64(b.x) && (g.x-16+g.velx) <= float64(b.x+b.w) {
				g.velx = 0
				b.color = color.RGBA{255, 255, 0, 255}
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(g.x, g.y)
	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

	moveX := false

	g.vely += 0.1

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.velx = 1
		moveX = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.velx = -1
		moveX = true
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.jump {
		g.vely = -3.5
		g.jump = true
	}
	if !moveX {
		g.velx = 0
	}

	for _, b := range g.rects {
		b.color = color.RGBA{255, 255, 255, 255}
	}

	g.collide(g.rects)

	g.x += g.velx
	g.y += g.vely

	g.tilemap.Draw(screen)

	if g.state != nil {
		g.state.Draw()
	}

	// for _, b := range g.rects {
	// 	b.Draw(screen)
	// }
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
