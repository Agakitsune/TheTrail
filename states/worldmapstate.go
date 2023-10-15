package states

import (
	"TheTrail/engine"
	"fmt"
	
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	
	"image"
	"image/color"
)

type WorldMapState struct {
	game          *engine.Game
	animator      *engine.Animator
	water         engine.MultiSprite
	reunion       engine.MultiSprite
	mapFlag       engine.MultiSprite
	flagAnimator  *engine.Animator
	flagPositions []engine.Vector2
	
	texts         []string

	dogica font.Face

	selectedFlag  int
	selector *ebiten.Image
}

func (s *WorldMapState) Load(g *engine.Game) {
	fmt.Println("WorldMapState Load")

	s.dogica = engine.LoadFont("./assets/fonts/dogica.ttf", 8)
	s.dogica = text.FaceWithLineHeight(s.dogica, 10)

	s.game = g

	// Water
	s.water = engine.MultiSprite{
		Sprites: []*ebiten.Image{
			engine.LoadImage("./assets/water_animated.png"),
		},
		Rect: image.Rect(0, 0, 320, 128),
		X:    0,
		Y:    0,
		Velx: 0,
		Vely: 0,
	}

	// Reunion
	s.reunion = engine.MultiSprite{
		Sprites: []*ebiten.Image{
			engine.LoadImage("./assets/caillou974.png"),
		},
		Rect: image.Rect(0, 0, 230, 204),
		X:    0,
		Y:    0,
		Velx: 0,
		Vely: 0,
	}
	s.reunion.X = 90

	s.selector = engine.LoadImage("./assets/selector.png")

	// Animator
	s.animator = &engine.Animator{
		Animations: map[string]*engine.Animation{
			"idle": &engine.Animation{
				Frames:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				Row:       0,
				LoopOn:    0,
				Selection: 0,
				Speed:     10,
			},
		},
		Current: "idle",
	}
	s.animator.SetFrameSize(320, 128)

	// Flags
	s.mapFlag = engine.MultiSprite{
		Sprites: []*ebiten.Image{
			engine.LoadImage("./assets/map_flag.png"),
		},
		Rect:        image.Rect(0, 0, 60, 60),
		X:           0,
		Y:           0,
		Velx:        0,
		Vely:        0,
		Scale:       engine.Vector2{0.3, 0.3},
		Initialized: true,
	}

	s.flagAnimator = &engine.Animator{
		Animations: map[string]*engine.Animation{
			"red": &engine.Animation{
				Frames:    []int{0, 1, 2, 3, 4},
				Row:       0,
				LoopOn:    0,
				Selection: 0,
				Speed:     10,
			},
			"white": &engine.Animation{
				Frames:    []int{5, 6, 7, 8, 9},
				Row:       0,
				LoopOn:    0,
				Selection: 0,
				Speed:     10,
			},
		},
		Current: "red",
	}
	s.flagAnimator.SetFrameSize(60, 60)

	s.texts = []string{
		"St Denis\n\nPoyox",
		"St André\n\nJe suis\nun tacocat",
		"St Untruc\n\nPtdr, t ki ?",
	}

	s.flagPositions = []engine.Vector2{
		engine.Vector2{X: 180, Y: 8},
		engine.Vector2{X: 220, Y: 20},
		engine.Vector2{X: 250, Y: 55},
	}
}

func (s *WorldMapState) Update() error {
	s.animator.Update(&s.water)
	s.flagAnimator.Update(&s.mapFlag)

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		s.selectedFlag++
		if s.selectedFlag >= len(s.flagPositions) {
			s.selectedFlag = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		s.selectedFlag--
		if s.selectedFlag < 0 {
			s.selectedFlag = len(s.flagPositions) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.game.SetState(&PlayState{})
	}

	return nil
}

func (s *WorldMapState) Draw(screen *ebiten.Image) {
	s.water.Y = 0
	s.water.Draw(screen)
	s.water.Y = 128
	s.water.Draw(screen)

	s.reunion.Draw(screen)

	for _, pos := range s.flagPositions {
		s.mapFlag.X = pos.X
		s.mapFlag.Y = pos.Y
		s.mapFlag.Draw(screen)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.4, 0.4)
	op.GeoM.Translate(s.flagPositions[s.selectedFlag].X - 8, s.flagPositions[s.selectedFlag].Y - 2)
	screen.DrawImage(s.selector, op);

	vector.DrawFilledRect(screen, 0, 0, engine.ScreenWidth/4, engine.ScreenHeight, color.RGBA{0, 0, 0, 200}, false)

	text.Draw(screen, s.texts[s.selectedFlag], s.dogica, 0, 8, color.RGBA{255, 255, 255, 255})
}
