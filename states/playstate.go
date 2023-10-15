package states

import (
	"TheTrail/engine"


	"github.com/hajimehoshi/ebiten/v2"
)

type PlayState struct {
	game *engine.Game
}

func (s *PlayState) Load(g *engine.Game) {
	s.game = g
	s.game.Init()
	// engine.NewColliderMap("./assets/collide.csv")
}

func (s *PlayState) Update() error {
	s.game.Animator.Update(&s.game.Dood)

	moveX := false
	run := true

	s.game.Dood.Vely += 0.1

	vel := 0.5
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		vel = 0.2
		run = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		s.game.Dood.Velx = vel
		s.game.Dood.Flip = false
		if !s.game.Airborne {
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
		if !s.game.Airborne {
			if run {
				s.game.Animator.SetAnimation("run")
			} else {
				s.game.Animator.SetAnimation("walk")
			}
		}
		moveX = true
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !s.game.Jump {
		s.game.Animator.SetAnimation("jump")
		s.game.Dood.Vely = -3
		s.game.Jump = true
		s.game.Airborne = true
	}
	if !moveX {
		if !s.game.Airborne {
			s.game.Animator.SetAnimation("idle")
		}
		s.game.Dood.Velx = 0
	}

	s.game.Collide(s.game.Collider.Boxes)

	s.game.Dood.X += s.game.Dood.Velx
	s.game.Dood.Y += s.game.Dood.Vely
	return nil
}

func (s *PlayState) Draw(screen *ebiten.Image) {
	s.game.Tilemap.Draw(screen)
	s.game.Dood.Draw(screen)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()), 0, 20)
}
