package engine

import (
	"os"
	"strconv"
	"strings"
)

type Rectangle struct {
	X int
	Y int
	Width int
	Height int
}

type Collider struct {
	Boxes []Rectangle
}

type lineInfo struct {
	x int
	y int
	width int
}

type rectInfo struct {
	x int
	y int
	width int
	height int
}

func NewColliderMap(path string)* Collider {
	var collider *Collider = new(Collider)
    data, err := os.ReadFile(path)

    if err != nil {
        panic(err)
    }

	found := false
	pos := 0
	width := 0

	info := make([]lineInfo, 0)

    lines := strings.Split(string(data), "\n")
    for i, line := range lines {
        if len(line) > 1 {
            for j, el := range strings.Split(line, ",") {
                val, err := strconv.Atoi(el)

				if err != nil {
					panic(err)
				}

                if val != -1 {
					if !found {
						pos = j
						found = true
					}
					width++
                } else {
					if found {
						info = append(info, lineInfo{pos, i, width})
						found = false
						width = 0
					}
				}
            }
			if found {
				info = append(info, lineInfo{pos, i, width})
			}
			found = false
			width = 0
        }
    }

	pending := make([]rectInfo, 0)
	for _, line := range info {
		newRect := true
		for i, rect := range pending {
			if rect.x == line.x && line.y == rect.y + rect.height {
				pending[i].height++
				newRect = false
				break
			}
		}

		if newRect {
			pending = append(pending, rectInfo{line.x, line.y, line.width, 1})
		}
	}

	for _, rect := range pending {
		collider.Boxes = append(collider.Boxes, Rectangle{rect.x * 8, rect.y * 8, rect.width * 8, rect.height * 8})
	}

	return collider
}
