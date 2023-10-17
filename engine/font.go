package engine

import (
	"io/ioutil"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func LoadFont(path string, size float64) font.Face {
	content, error := ioutil.ReadFile(path)
	if error != nil {
		log.Fatal(error)
	}

	tt, err := opentype.Parse(content)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	font, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	return font
}
