package states

import (
	"TheTrail/engine"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type EditorState struct {
	game *engine.Game

	heads []*ebiten.Image
    indexH int

    torso []*ebiten.Image
    indexT int

    boots []*ebiten.Image
    indexB int

    arrow *ebiten.Image

    arrowRectL image.Rectangle
        arrowRectR image.Rectangle

    }

    func (s *EditorState) Load(g *engine.Game) {
        s.game = g

        s.heads = []*ebiten.Image{
            engine.LoadImage("./assets/dood/head_one.png"),
            engine.LoadImage("./assets/dood/head_two.png"),
            engine.LoadImage("./assets/dood/head_three.png"),
        }
        s.indexH = 0

        s.torso = []*ebiten.Image {
            engine.LoadImage("./assets/dood/torso_one.png"),
            engine.LoadImage("./assets/dood/torso_two.png"),
            engine.LoadImage("./assets/dood/torso_three.png"),
        }
        s.indexT = 0

        s.boots = []*ebiten.Image {
            engine.LoadImage("./assets/dood/boots_one.png"),
            engine.LoadImage("./assets/dood/boots_two.png"),
            engine.LoadImage("./assets/dood/boots_three.png"),
        }
        s.indexB = 0

        s.arrow = engine.LoadImage("./assets/dood/arrow.png")

        s.arrowRectL = image.Rect(40, 0, 80, 40)
        s.arrowRectR = image.Rect(0, 0, 40, 40)
        fmt.Println("EditorState Load")
    }

    func (s *EditorState) Update() error {
        if ebiten.IsKeyPressed(ebiten.KeyA) {
            fmt.Println("EditorState Update")
            s.game.SetState(&PlayState{})
        }
        return nil;
    }

    func (s *EditorState) Draw(screen *ebiten.Image) {
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(0, 0)
        scalingFactor := 4.0
        op.GeoM.Scale(scalingFactor, scalingFactor)

        screen.DrawImage(s.heads[s.indexH].SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)
        screen.DrawImage(s.torso[s.indexT].SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)
        screen.DrawImage(s.boots[s.indexB].SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)
        op.GeoM.Reset()

        for i := 0; i < 3; i++ {
            op.GeoM.Translate(35, float64(100 + 95 * i))
            scalingFactorArrowL := 0.4
            op.GeoM.Scale(scalingFactorArrowL, scalingFactorArrowL)
            screen.DrawImage(ebiten.NewImageFromImage(s.arrow.SubImage(s.arrowRectL)), op)
            op.GeoM.Reset()
        }

        for i :=0; i < 3; i++ {
            scalingFactorArrowR := 0.4
            op.GeoM.Translate(250, float64(100 + 95 * i))
            op.GeoM.Scale(scalingFactorArrowR, scalingFactorArrowR)
            screen.DrawImage(ebiten.NewImageFromImage(s.arrow.SubImage(s.arrowRectR)), op)
            op.GeoM.Reset()
        }


// 	fmt.Println("EditorState Draw")

}
