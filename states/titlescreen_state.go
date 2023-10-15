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
	s.image = engine.LoadImage("./assets/back_langevin.png")
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
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		fmt.Println("TitleScreenState Update")
		s.game.SetState(&EditorState{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Println("you just pressed a Enter key")

		fmt.Println("PLAY FADE INTO CHARACTER EDITOR")
		character_edit = true

	}
	return nil
}

func (s *TitleScreenState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{135, 206, 235, 0xff})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.33, 0.27)

	screen.DrawImage(s.image, op)
	// vector.DrawFilledRect(screen, 80, 10, 160, 25, color.Black, false)

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
			// x := 50
			text.Draw(screen, l, titleArcadeFont, x, (i+2)*titleFontSize, color.Black)
		}
		for i, l := range texts {
			// x := (engine.ScreenHeight - len(l)*fontSize) / 2
			x := 60
			text.Draw(screen, l, arcadeFont, x, (i+12)*fontSize, color.White)
		}
	}

	// fmt.Println("TitleScreenState Draw")
}
