package states

import (
	"TheTrail/engine"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	camera "github.com/melonfunction/ebiten-camera"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	fontSize      = 10
	titleFontSize = fontSize * 1.5
)

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
	character_edit  = false
)

type TitleScreenState struct {
	game *engine.Game

	image *ebiten.Image
}

func (s *TitleScreenState) Load(g *engine.Game) {
	s.image = engine.LoadImage("./assets/back_volcan.png")
	s.game = g

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	fmt.Println("TitleScreenState Load")
}

func (s *TitleScreenState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Println("TitleScreenState Update")
		s.game.SetState(&EditorState{})
		character_edit = true
	}
	return nil
}

func (s *TitleScreenState) Draw(screen *ebiten.Image, camera *camera.Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.05, 0.06)

	screen.DrawImage(s.image, op)
	if character_edit == false {
		vector.DrawFilledRect(screen, 80, 10, 160, 25, color.RGBA{189, 195, 199, 255}, false)
	}
	//load picture assets\back_langevin.png for wallpaper

	var titleTexts []string
	var texts []string

	titleTexts = []string{"The Trail"}
	texts = []string{"Press ENTER to play"}

	if character_edit == false {
		for i, l := range titleTexts {
			x := (engine.ScreenWidth - len(l)*titleFontSize) / 2
			text.Draw(screen, l, titleArcadeFont, x, (i+2)*titleFontSize, color.Black)
		}

		// draw character walking here
		// HERE

		for i, l := range texts {
			x := 60
			text.Draw(screen, l, arcadeFont, x, (i+16)*fontSize, color.White)
		}
	}

}
