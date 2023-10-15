package engine

import (
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tilemap struct {
	tileset *ebiten.Image

	tiles []int

	Width  int32
	Height int32
}

func NewTilemap(mapPath, tilesetPath string) *Tilemap {
	var gmap *Tilemap = new(Tilemap)
	data, err := os.ReadFile(mapPath)

	if err != nil {
		panic(err)
	}

	gmap.tiles = make([]int, 0)
	gmap.tileset = LoadImage(tilesetPath)
	lines := strings.Split(strings.Replace(string(data), "\r", "", -1), "\n")
	for i, line := range lines {
		if len(line) > 1 {
			for _, el := range strings.Split(line, ",") {
				val, err := strconv.Atoi(el)
				if err == nil {
					gmap.tiles = append(gmap.tiles, val)
				}
				if i == 0 {
					gmap.Width++
				}
			}
			gmap.Height++
		}
	}
	return gmap
}

func (self *Tilemap) Draw(screen *ebiten.Image) {
	for i, tile := range self.tiles {
		if tile != -1 {
			op := &ebiten.DrawImageOptions{}
			sx, sy := ((tile & 7) * 8), (tile>>3)*8
			op.GeoM.Translate(float64(i%40)*8, float64(i/40)*8)
			screen.DrawImage(self.tileset.SubImage(image.Rect(sx, sy, sx+8, sy+8)).(*ebiten.Image), op)
		}
	}
}
