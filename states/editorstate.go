package states

import (
	"TheTrail/engine"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	camera "github.com/melonfunction/ebiten-camera"
)

type EditorState struct {
	game *engine.Game

	currentChoice int

	heads  []*ebiten.Image
	indexH int

	torso  []*ebiten.Image
	indexT int

	boots  []*ebiten.Image
	indexB int

	arrow *ebiten.Image

	arrowRectL        image.Rectangle
	arrowRectR        image.Rectangle
	disableArrowRectL image.Rectangle
	disableArrowRectR image.Rectangle

	title   *ebiten.Image
	goToMap *ebiten.Image
}

func (s *EditorState) Load(g *engine.Game) {
	s.game = g

	s.currentChoice = 0
	s.heads = []*ebiten.Image{
		engine.LoadImage("./assets/dood/head_one.png"),
		engine.LoadImage("./assets/dood/head_two.png"),
		engine.LoadImage("./assets/dood/head_three.png"),
	}
	s.indexH = 0

	s.torso = []*ebiten.Image{
		engine.LoadImage("./assets/dood/torso_one.png"),
		engine.LoadImage("./assets/dood/torso_two.png"),
		engine.LoadImage("./assets/dood/torso_three.png"),
	}
	s.indexT = 0

	s.boots = []*ebiten.Image{
		engine.LoadImage("./assets/dood/boots_one.png"),
		engine.LoadImage("./assets/dood/boots_two.png"),
		engine.LoadImage("./assets/dood/boots_three.png"),
	}
	s.indexB = 0

	s.arrow = engine.LoadImage("./assets/dood/arrow.png")

	s.arrowRectL = image.Rect(40, 0, 80, 40)
	s.arrowRectR = image.Rect(0, 0, 40, 40)

	s.disableArrowRectL = image.Rect(40, 40, 80, 80)
	s.disableArrowRectR = image.Rect(0, 40, 40, 80)

	s.title = engine.LoadImage("./assets/dood/title.png")
	s.goToMap = engine.LoadImage("./assets/dood/goToMap.png")
	fmt.Println("EditorState Load")
}

func (s *EditorState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Println("EditorState Update")
		s.game.SetState(&WorldMapState{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if s.currentChoice == 0 {
			s.currentChoice = 2
		} else {
			s.currentChoice--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if s.currentChoice == 2 {
			s.currentChoice = 0
		} else {
			s.currentChoice++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if s.currentChoice == 0 {
			if s.indexH == 0 {
				s.indexH = 2
			} else {
				s.indexH--
			}
		}
		if s.currentChoice == 1 {
			if s.indexT == 0 {
				s.indexT = 2
			} else {
				s.indexT--
			}
		}
		if s.currentChoice == 2 {
			if s.indexB == 0 {
				s.indexB = 2
			} else {
				s.indexB--
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if s.currentChoice == 0 {
			if s.indexH == 2 {
				s.indexH = 0
			} else {
				s.indexH++
			}
		}
		if s.currentChoice == 1 {
			if s.indexT == 2 {
				s.indexT = 0
			} else {
				s.indexT++
			}
		}
		if s.currentChoice == 2 {
			if s.indexB == 2 {
				s.indexB = 0
			} else {
				s.indexB++
			}
		}
	}
	return nil
}

func (s *EditorState) Draw(screen *ebiten.Image, camera *camera.Camera) {

	screen.Fill(color.RGBA{95, 181, 172, 100})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(engine.ScreenWidth/1.05, engine.ScreenHeight/1.03)
	op.GeoM.Scale(0.45, 0.45)
	screen.DrawImage(s.goToMap, op)
	op.GeoM.Reset()

	op.GeoM.Translate(engine.ScreenWidth/4, 3)
	screen.DrawImage(s.title, op)
	op.GeoM.Reset()

	op.GeoM.Translate(2, 8)
	scalingFactor := 4.0
	op.GeoM.Scale(scalingFactor, scalingFactor)

	screen.DrawImage(s.heads[s.indexH].SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)
	screen.DrawImage(s.torso[s.indexT].SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)
	screen.DrawImage(s.boots[s.indexB].SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)
	op.GeoM.Reset()

	for i := 0; i < 3; i++ {
		op.GeoM.Translate(40, float64(175+95*i))
		scalingFactorArrowL := 0.4
		op.GeoM.Scale(scalingFactorArrowL, scalingFactorArrowL)
		if s.currentChoice == i {
			screen.DrawImage(ebiten.NewImageFromImage(s.arrow.SubImage(s.arrowRectL)), op)
		} else {
			screen.DrawImage(ebiten.NewImageFromImage(s.arrow.SubImage(s.disableArrowRectL)), op)
		}
		op.GeoM.Reset()
	}

	for i := 0; i < 3; i++ {
		scalingFactorArrowR := 0.4
		op.GeoM.Translate(280, float64(175+95*i))
		op.GeoM.Scale(scalingFactorArrowR, scalingFactorArrowR)
		if s.currentChoice == i {
			screen.DrawImage(ebiten.NewImageFromImage(s.arrow.SubImage(s.arrowRectR)), op)
		} else {
			screen.DrawImage(ebiten.NewImageFromImage(s.arrow.SubImage(s.disableArrowRectR)), op)
		}
		op.GeoM.Reset()
	}
}
