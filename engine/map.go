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

	tiles []*ebiten.Image

	Width int
  Height int
}

func NewTilemap(mapPath, tilesetPath string) *Tilemap {
	var gmap *Tilemap = new(Tilemap)
	data, err := os.ReadFile(mapPath)

	if err != nil {
		panic(err)
	}

    tilesIndex := make([]int, 0)
    gmap.tiles = make([]*ebiten.Image, 0)
	  gmap.tileset = LoadImage(tilesetPath)
    lines := strings.Split(strings.Replace(string(data), "\r", "", -1), "\n")
    for i, line := range lines {
        if len(line) > 1 {
            for _, el := range strings.Split(line, ",") {
                val, err := strconv.Atoi(el)
                if err == nil {
                    tilesIndex = append(tilesIndex, val)
                }
                if i == 0 {
                    gmap.Width++
                }
            }
            gmap.Height++
        }
    }

    for _, tile := range tilesIndex {
		if tile != -1 {
			sx, sy := ((tile % 12) * 8), (tile / 12) * 8
            gmap.tiles = append(gmap.tiles, gmap.tileset.SubImage(image.Rect(sx, sy, sx+8, sy+8)).(*ebiten.Image))
		} else {
            gmap.tiles = append(gmap.tiles, nil)
        }
	}
    
    return gmap
}

func (self *Tilemap) Draw(screen *ebiten.Image) {
    
    for i, tile := range self.tiles {
        if tile == nil {
            continue
        }
        x := (i % self.Width) * 8
        y := (i / self.Width) * 8
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(float64(x), float64(y))
        screen.DrawImage(tile, op)
    }
}
