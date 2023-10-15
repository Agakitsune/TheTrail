package states

import (
	"fmt"

	"TheTrail/engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	camera "github.com/melonfunction/ebiten-camera"

	// raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
)

type PlayState struct {
	game		*engine.Game
	music 	  	*engine.Audio
}

func (s *PlayState) Load(gm *engine.Game) {
	s.game = gm
	s.game.Init()
	s.music = engine.CreateAudio("assets/EpitechGameJam-_In_Game.ogg")
	s.music.Play()
}

func (s *PlayState) Update() error {
	s.game.Animator.Update(s.game.Dood)

	moveX := false
	run := true

	if s.game.Dood.Vely < 4.0 {
		s.game.Dood.Vely += 0.1
	}

	if s.game.SceneTransition {
		x := (1-s.game.Timer)*float64(s.game.SceneX) + s.game.Timer*float64(s.game.ToSceneX)
		y := (1-s.game.Timer)*float64(s.game.SceneY) + s.game.Timer*float64(s.game.ToSceneY)

		s.game.Timer += 0.05

		s.game.Cam.SetPosition(160+x, 90+y)

		if s.game.Timer >= 1 {
			s.game.SceneX = s.game.ToSceneX
			s.game.SceneY = s.game.ToSceneY
			s.game.Cam.SetPosition(160.0+float64(s.game.SceneX), 90.0+float64(s.game.SceneY))
			s.game.SceneTransition = false
		}
	} else {

		vel := 1.0
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			vel = 0.2
			run = false
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			s.game.Dood.Velx = vel
			s.game.Dood.Flip = false
			if !s.game.Dood.Airborne {
				if run {
					s.game.Animator.SetAnimation("run")
				} else {
					s.game.Animator.SetAnimation("walk")
				}
			}
			moveX = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			s.game.Dood.Velx = -vel
			s.game.Dood.Flip = true
			if !s.game.Dood.Airborne {
				if run {
					s.game.Animator.SetAnimation("run")
				} else {
					s.game.Animator.SetAnimation("walk")
				}
			}
			moveX = true
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !s.game.Dood.Jump {
			s.game.Animator.SetAnimation("jump")
			s.game.Dood.Vely = -3
			s.game.Dood.Jump = true
			s.game.Dood.Airborne = true
		}
		if !moveX {
			if !s.game.Dood.Airborne {
				s.game.Animator.SetAnimation("idle")
			}
			s.game.Dood.Velx = 0
		}

	}

	for _, collider := range s.game.Collider {
		collider.Update(s.game, s.game.Dood)
	}

	// s.game.Collider[s.game.SceneX].Update(s.game, s.game.Dood)
	s.game.Scene.Update(s.game, s.game.Dood)

	if s.game.SceneTransition {
		if s.game.SceneY >= s.game.ToSceneY {
			s.game.Dood.Y += s.game.Dood.Vely
		}
		return nil
	}

	s.game.Dood.X += s.game.Dood.Velx
	s.game.Dood.Y += s.game.Dood.Vely

	return nil
}

func (s *PlayState) Draw(screen *ebiten.Image, camera *camera.Camera) {
	for _, tilemap := range s.game.Tilemap {
		tilemap.Draw(screen, camera)
	}

	// s.game.Tilemap[s.game.SceneX].Draw(screen, camera)
	s.game.Dood.Draw(screen, camera)

	if s.game.Debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()), 0, 20)
		for _, collider := range s.game.Collider {
			collider.Draw(screen)
		}
	}

}
