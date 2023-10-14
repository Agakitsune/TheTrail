package engine

import (
	"bytes"
	"io/ioutil"
	"log"
	
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(path string) *ebiten.Image {
	content, error := ioutil.ReadFile(path)
	if error != nil {
		log.Fatal(error)
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}
