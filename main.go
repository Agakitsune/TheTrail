package main

import (
	_ "image/png"
	"log"

	"TheTrail/engine"
	"TheTrail/states"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var state engine.State = &states.PlayState{}
	var game = engine.CreateGame(state)

	ebiten.SetWindowSize(engine.ScreenWidth*2, engine.ScreenHeight*2)
	ebiten.SetWindowTitle("The Trail")
	ebiten.SetMaxTPS(90)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
