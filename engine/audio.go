package engine

import (
	"bufio"
	"io"
	"log"
	"os"
	"bytes"

	// "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	sampleRate     = 48000
	bytesPerSample = 4 // 2 channels * 2 bytes (16 bit)

	introLengthInSecond = 0
	loopLengthInSecond  = 72
)

type Audio struct {
	player       *audio.Player
	audioContext *audio.Context
	rawMusicData []byte
}

func (aud *Audio) Play() {
	// Play the infinite-length stream. This never endaudio.
	aud.player.Play()
}

func (aud *Audio) Pause() {
	aud.player.Pause()
}

func CreateAudio(path string) *Audio {
	aud := &Audio{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		log.Fatal(err)
		return nil
	}

	aud.rawMusicData = bs

	file.Close()

	// Play musik
	if aud.audioContext == nil {
		aud.audioContext = audio.NewContext(sampleRate)
	}

	// Decode an Ogg file.
	// oggS is a decoded io.ReadCloser and io.Seeker.
	oggS, err := vorbis.DecodeWithoutResampling(bytes.NewReader(aud.rawMusicData))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Create an infinite loop stream from the decoded byteaudio.
	// s is still an io.ReadCloser and io.Seeker.
	stream := audio.NewInfiniteLoopWithIntro(oggS, introLengthInSecond*bytesPerSample*sampleRate, loopLengthInSecond*bytesPerSample*sampleRate)

	aud.player, err = aud.audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return aud
}
