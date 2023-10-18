package states

import (
	"TheTrail/engine"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	camera "github.com/melonfunction/ebiten-camera"
)

type PlayState struct {
	game  *engine.Game
	music *engine.Audio

	deathTransition bool
	deathTimer int
	deathBlack *ebiten.Image

	win bool
	winTimer int
	artTimer int
	winArt *ebiten.Image
}

func (s *PlayState) Load(g *engine.Game) {
	s.game = g

	s.game.Collider = make([]*engine.Collider, 0)
	s.game.Collider = append(s.game.Collider, engine.NewColliderMap("./level/level1_collide.csv", 0, -548))
	// s.game.Collider = append(s.game.Collider, engine.NewColliderMap("./assets/next_collide.csv", 320, 0))

	s.game.Tilemap = make([]*engine.Tilemap, 0)
	s.game.Tilemap = append(s.game.Tilemap, engine.NewTilemap("./level/level1_back.csv", "./assets/new_grass.png", 0, -548))
	s.game.Tilemap = append(s.game.Tilemap, engine.NewTilemap("./level/level1_ground.csv", "./assets/new_grass.png", 0, -548))

	s.game.Front = make([]*engine.Tilemap, 0)
	s.game.Front = append(s.game.Front, engine.NewTilemap("./level/level1_front.csv", "./assets/new_grass.png", 0, -548))
	// s.game.Tilemap = append(s.game.Tilemap, engine.NewTilemap("./assets/next_draw.csv", "./assets/new_grass.png", 320, 0))

	s.deathBlack = engine.LoadImage("./assets/black.png")

	s.winArt = engine.LoadImage("./assets/art.png");

	// s.game.Dood.X = 8
	// s.game.Dood.Y = 87 * 8

	// s.game.Init()
	s.music = engine.CreateAudio("assets/EpitechGameJam-_In_Game.ogg")
	s.music.Play()
	// engine.NewColliderMap("./assets/collide.csv")
}

func (s *PlayState) Update() error {
	s.game.Animator.Update(s.game.Dood)

	moveX := false
	moveY := false
	run := true

	if !s.game.Dood.Climbing && !s.deathTransition {
		if s.game.Dood.SlowFall {
			if s.game.Dood.Vely < 0.3 {
				s.game.Dood.Vely += 0.05
			}
		} else {
			if s.game.Dood.Vely < 4.0 {
				s.game.Dood.Vely += 0.1
			}
		}
	}

	if s.win {
		s.winTimer++
		if s.winTimer < 10 {
			s.game.Dood.Velx = 0
			s.game.Dood.Flip = false
		} else if s.winTimer < 180 {
			s.game.Animator.SetAnimation("walk")
			s.game.Dood.Velx = 0.2
		} else {
			s.game.Animator.SetAnimation("idle")
			s.game.Dood.Velx = 0
			if s.winTimer >= 240 {
				// if s.winTimer % 10 == 0 {
					s.artTimer++
				// }
			}
		}
	} else if s.deathTransition {
		s.deathTimer++
		if s.deathTimer == 30 {
			s.game.Cam.SetPosition(160, 90)
			s.game.SceneX = 0
			s.game.SceneY = 0
			s.game.Dood.Dead = false
			s.game.Animator.SetAnimation("spawn")
			s.game.Dood.X = 32
			s.game.Dood.Y = 128
			s.game.Scene.ShiftX = 0
			s.game.Scene.ShiftY = 0
		} else if s.deathTimer >= 60 {
			if (s.game.Animator.OnLastFrame()) {
				s.game.Animator.SetAnimation("idle")
				s.deathTransition = false
				s.deathTimer = 0
			}
		}
	} else if s.game.SceneTransition {
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

		if s.game.Dood.Win {
			s.win = true
		} else if s.game.Dood.Dead {
			s.game.Dood.Velx = 0
			s.game.Dood.Vely = 0

			s.deathTransition = true

			s.game.Animator.SetAnimation("death")

		} else if !s.game.Dood.Climbing {
			vel := 1.0

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

			if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
				s.game.Dood.Falltrough = true
			} else {
				s.game.Dood.Falltrough = false
			}

			if ebiten.IsKeyPressed(ebiten.KeySpace) && !s.game.Dood.Jump {
				s.game.Animator.SetAnimation("jump")
				s.game.Dood.Vely = -3
				s.game.Dood.Jump = true
				s.game.Dood.Airborne = true
			}

			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				s.game.Dood.TryClimb = s.game.Dood.Stamina > 0
			} else {
				s.game.Dood.TryClimb = false
			}

			if s.game.Dood.SlowFall {
				s.game.Animator.SetAnimation("climb")
				s.game.Dood.SlowFall = s.game.Dood.Dir == int(s.game.Dood.Velx)
			}

			if !moveX {
				if !s.game.Dood.Airborne {
					s.game.Animator.SetAnimation("idle")
				}
				s.game.Dood.Velx = 0
				s.game.Dood.Dir = 0
			}
		} else {

			s.game.Dood.Climbing = ebiten.IsKeyPressed(ebiten.KeyShift)

			s.game.Dood.Stamina--

			println(s.game.Dood.Stamina)

			if !s.game.Dood.Edge && (ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp)) {
				s.game.Dood.Vely = -0.2
				moveY = true
			}

			if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
				s.game.Dood.Vely = 0.2
				moveY = true
			}

			if ebiten.IsKeyPressed(ebiten.KeySpace) && !s.game.Dood.Jump {
				// s.game.Animator.SetAnimation("jump")
				if s.game.Dood.Flip && ebiten.IsKeyPressed(ebiten.KeyD) && s.game.Dood.Stamina >= 10 {
					s.game.Dood.Vely = -2
					s.game.Dood.Velx = 2
					s.game.Dood.Stamina -= 10
				} else if !s.game.Dood.Flip && ebiten.IsKeyPressed(ebiten.KeyA) && s.game.Dood.Stamina >= 10 {
					s.game.Dood.Vely = -2
					s.game.Dood.Velx = -2
					s.game.Dood.Stamina -= 10
				}
				s.game.Dood.Climbing = false
				s.game.Dood.Jump = true
				s.game.Dood.Airborne = true
				moveY = true
			}

			if s.game.Dood.Stamina <= 0 {
				s.game.Dood.Climbing = false
				s.game.Dood.SlowFall = false
			}

			if !moveY {
				s.game.Dood.Vely = 0
				if s.game.Dood.Stamina <= 50 {
					s.game.Dood.Vely = 0.1
				}
			}

		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) && !s.game.Dood.Dead {
		s.game.Dood.Dead = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) && !s.game.Dood.Win {
		s.game.Dood.X = 1440
		s.game.Dood.Y = -270
		s.game.Cam.SetPosition(1440, -270)
		s.game.Scene.ShiftX = 4
		s.game.Scene.ShiftY = -2
	}

	for _, collider := range s.game.Collider {
		collider.Update(s.game, s.game.Dood)
	}

	if s.game.Dood.Climbing {
		s.game.Animator.SetAnimation("climb")
	}

	// s.game.Collider[s.game.SceneX].Update(s.game, s.game.Dood)
	if !s.deathTransition {
		s.game.Scene.Update(s.game, s.game.Dood)
	}

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
	op := &ebiten.DrawImageOptions{}
	op = camera.GetTranslation(op, s.game.Cam.X - 160, -360)
	// op.GeoM.Translate(0, float64(-s.game.SceneY * 180) - 360)
	screen.DrawImage(s.game.Background, op)
	
	for _, tilemap := range s.game.Tilemap {
		tilemap.Draw(screen, camera)
	}
	
	s.game.Dood.Draw(screen, camera)

	for _, tilemap := range s.game.Front {
		tilemap.Draw(screen, camera)
	}

	if s.game.Debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()), 0, 20)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Camera: %0.2f, %0.2f", s.game.Cam.X, s.game.Cam.Y), 0, 40)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pos: %0.2f, %0.2f", s.game.Dood.X, s.game.Dood.Y), 0, 60)
		for _, collider := range s.game.Collider {
			collider.Draw(screen, camera)
		}
	}

	if s.deathTransition {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(-960 + s.deathTimer * 22), 0)
		screen.DrawImage(s.deathBlack, op)
	}

	if s.win && s.artTimer > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-106, -71)
		if s.artTimer < 120 {
			value := -0.01 * float64(s.artTimer) + 2.1
			op.GeoM.Scale(float64(value), float64(value))
			op.GeoM.Rotate(-float64(s.artTimer) / 2 * (3.14 / 180) + (40 * (3.14 / 180)))
		} else {
			op.GeoM.Scale(0.9, 0.9)
			op.GeoM.Rotate(-20 * (3.14 / 180))
		}
		op.GeoM.Translate(320/2, 180/2)
		screen.DrawImage(s.winArt, op)
		if (s.artTimer >= 180) {
			s.game.SetState(&WorldMapState{})
		}
	}

}
