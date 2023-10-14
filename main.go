package main

import (
	_ "image/png"
	"log"

	"TheTrail/engine"
	"TheTrail/states"

	// "fmt"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	var state engine.State = &states.PlayState{}
	var game = engine.CreateGame(state)

	ebiten.SetWindowSize(engine.ScreenWidth*2, engine.ScreenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
