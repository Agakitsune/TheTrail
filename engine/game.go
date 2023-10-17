package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

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
	Dood* MultiSprite

	Collider []*Collider
	Tilemap []*Tilemap

	Scene *SceneTrigger

	SceneTransition bool
	SceneX int
	SceneY int
	ToSceneX int
	ToSceneY int
	Timer float64
 
	Animator *Animator

	Cam *camera.Camera

	Debug bool

	state State
}

func (g *Game) SetState(state State) {
	g.state = state
	g.state.Load(g)
}

func (g *Game) Init() {
	g.Collider = make([]*Collider, 0)
	g.Collider = append(g.Collider, NewColliderMap("./assets/new_collision.csv", 0, 0))
	g.Collider = append(g.Collider, NewColliderMap("./assets/next_collide.csv", 320, 0))

	g.Tilemap = make([]*Tilemap, 0)
	g.Tilemap = append(g.Tilemap, NewTilemap("./assets/new_draw.csv", "./assets/grass.png", 0, 0))
	g.Tilemap = append(g.Tilemap, NewTilemap("./assets/next_draw.csv", "./assets/tileset.png", 320, 0))

	g.Scene = NewSceneTrigger(320, 180)

	g.Cam = camera.NewCamera(320, 180, 160, 90, 0, 1)

	g.Dood = NewSprite(
		[]string{
			"./assets/dood/boots_one.png",
			"./assets/dood/torso_three.png",
			"./assets/dood/head_two.png",
		},
		image.Rect(0, 0, 32, 32),
		Vector2{1, 1},
	)

	g.Dood.X = 32
	g.Dood.Y = 32

	g.Dood.Jump = false;
	g.Dood.Airborne = false;

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
			"climb": &Animation{
				Frames:    []int{0},
				Row:       4,
				LoopOn:    0,
				Selection: 0,
				Speed:     1,
			},
		},
		Current: "idle",
	}

	g.Animator.SetFrameSize(32, 32)

	// Initialize the state
	// g.state = nil
}

func (g *Game) MakeSceneTransition(x, y int) {
	g.SceneTransition = true
	g.ToSceneX = x
	g.ToSceneY = y
	g.Timer = 0
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.Debug = !g.Debug
	}

	if g.state != nil {
		return g.state.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state != nil {
		g.Cam.Surface.Clear()
		g.state.Draw(g.Cam.Surface, g.Cam)

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
