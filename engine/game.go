package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	camera "github.com/melonfunction/ebiten-camera"
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

type Game struct {
	Dood MultiSprite

	Collider *Collider

	Jump     bool
	Airborne bool

	Tilemap  *Tilemap
	Animator *Animator

	Cam *camera.Camera

	state State
}

func (g *Game) SetState(state State) {
	g.state = state
	g.state.Load(g)
}

func (g *Game) Init() {
	g.Collider = NewColliderMap("./assets/collide.csv")
	g.Jump = false

	g.Cam = camera.NewCamera(320, 180, 0, 0, 0, 1)

	g.Dood = MultiSprite{
		Sprites: []*ebiten.Image{
			LoadImage("./assets/dood/boots_one.png"),
			LoadImage("./assets/dood/torso_three.png"),
			LoadImage("./assets/dood/head_two.png"),
		},
		Rect: image.Rect(0, 0, 32, 32),
		X:    0,
		Y:    0,
		Velx: 0,
		Vely: 0,
	}

	g.Animator = &Animator{
		Animations: map[string]*Animation{
			"idle": &Animation{
				Frames:    []int{0, 1, 2, 3, 4, 5},
				Row:       0,
				LoopOn:    0,
				Selection: 0,
				Speed:     10,
			},
			"walk": &Animation{
				Frames:    []int{0, 1, 2, 3, 4, 5, 6},
				Row:       1,
				LoopOn:    1,
				Selection: 0,
				Speed:     10,
			},
			"run": &Animation{
				Frames:    []int{0, 1, 2, 3, 4, 5},
				Row:       2,
				LoopOn:    0,
				Selection: 0,
				Speed:     10,
			},
			"jump": &Animation{
				Frames:    []int{0, 1, 2},
				Row:       3,
				LoopOn:    2,
				Selection: 0,
				Speed:     5,
			},
		},
		Current: "idle",
	}

	g.Animator.SetFrameSize(32, 32)

	g.Tilemap = NewTilemap("./assets/map.csv", "./assets/tileset.png")

	// Initialize the state
	// g.state = nil
}

func (this *Game) Collide(boxes []Rectangle) {
	for _, b := range boxes {
		if this.Dood.Vely != 0 {
			// println("velx: ", fmt.Sprintf("%f", this.Dood.Velx))
			// println("gx: ", fmt.Sprintf("%f", this.Dood.X))
			// println("bx: ", fmt.Sprintf("%f", b.X))
			// println("bx: ", fmt.Sprintf("%f", b.X + b.Width))
			// println("estimateMin: ", fmt.Sprintf("%f", (this.Dood.X + 32 + this.Dood.Velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (this.Dood.X - 32 + this.Dood.Velx)))
			if (this.Dood.X+22+this.Dood.Velx) <= float64(b.X) || (this.Dood.X+11+this.Dood.Velx) >= float64(b.X+b.Width) {
				continue
			}
			// println("vely: ", fmt.Sprintf("%f", this.Dood.Vely))
			// println("gy: ", fmt.Sprintf("%f", this.Dood.Y))
			// println("by: ", fmt.Sprintf("%f", b.Y))
			// println("by + bh: ", fmt.Sprintf("%f",b.Y + b.Height))
			// println("estimateMin: ", fmt.Sprintf("%f", (this.Dood.Y + 32 + this.Dood.Vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (this.Dood.Y - 32 + this.Dood.Vely)))
			if (this.Dood.Y+32+this.Dood.Vely) >= float64(b.Y) && (this.Dood.Y+7+this.Dood.Vely) <= float64(b.Y+b.Height) {
				if this.Dood.Vely > 0 {
					this.Jump = false
					this.Airborne = false
				}
				this.Dood.Vely = 0
			}
		}
		if this.Dood.Velx != 0 {
			// println("vely: ", fmt.Sprintf("%f", this.Dood.Vely))
			// println("gy: ", fmt.Sprintf("%f", this.Dood.Y))
			// println("by: ", fmt.Sprintf("%f", b.Y))
			// println("by + bh: ", fmt.Sprintf("%f",b.Y + b.Height))
			// println("estimateMin: ", fmt.Sprintf("%f", (this.Dood.Y + 32 + this.Dood.Vely)))
			// println("estimateMax: ", fmt.Sprintf("%f", (this.Dood.Y - 32 + this.Dood.Vely)))
			if (this.Dood.Y+32+this.Dood.Vely) <= float64(b.Y) || (this.Dood.Y+7+this.Dood.Vely) >= float64(b.Y+b.Height) {
				continue
			}
			// println("velx: ", fmt.Sprintf("%f", this.Dood.Velx))
			// println("gx: ", fmt.Sprintf("%f", this.Dood.X))
			// println("bx: ", fmt.Sprintf("%f", b.X))
			// println("bx: ", fmt.Sprintf("%f", b.X + b.Width))
			// println("estimateMin: ", fmt.Sprintf("%f", (this.Dood.X + 32 + this.Dood.Velx)))
			// println("estimateMax: ", fmt.Sprintf("%f", (this.Dood.X - 32 + this.Dood.Velx)))
			if (this.Dood.X+22+this.Dood.Velx) >= float64(b.X) && (this.Dood.X+11+this.Dood.Velx) <= float64(b.X+b.Width) {
				this.Dood.Velx = 0
			}
		}
	}
}

func (g *Game) Update() error {
	if g.state != nil {
		return g.state.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state != nil {
		g.Cam.Surface.Clear()
		g.state.Draw(g.Cam.Surface)

		g.Cam.Blit(screen)
	}
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
