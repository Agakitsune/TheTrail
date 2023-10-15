package states

import (
	"TheTrail/engine"

	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	camera "github.com/melonfunction/ebiten-camera"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	// raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
)

const (
	screenWidth    = 640
	screenHeight   = 480
	sampleRate     = 48000
	bytesPerSample = 4 // 2 channels * 2 bytes (16 bit)

	introLengthInSecond = 0
	loopLengthInSecond  = 72
)

type PlayState struct {
	game         *engine.Game
	player       *audio.Player
	audioContext *audio.Context
	rawMusicData []byte
}

func (s *PlayState) Load(gm *engine.Game) {
	s.game = gm
	s.game.Init()

	file, err := os.Open("assets/EpitechGameJam-_In_Game.ogg")
	if err != nil {
		log.Fatal(err)
	}

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	s.rawMusicData = bs

	file.Close()

	// Play musik
	if s.audioContext == nil {
		s.audioContext = audio.NewContext(sampleRate)
	}

	// Decode an Ogg file.
	// oggS is a decoded io.ReadCloser and io.Seeker.
	oggS, err := vorbis.DecodeWithoutResampling(bytes.NewReader(s.rawMusicData))
	if err != nil {
		panic(err)
	}

	// Create an infinite loop stream from the decoded bytes.
	// s is still an io.ReadCloser and io.Seeker.
	stream := audio.NewInfiniteLoopWithIntro(oggS, introLengthInSecond*bytesPerSample*sampleRate, loopLengthInSecond*bytesPerSample*sampleRate)

	s.player, err = s.audioContext.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	// Play the infinite-length stream. This never ends.
	s.player.Play()
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
