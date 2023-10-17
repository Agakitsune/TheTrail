package engine

type SceneTrigger struct {
	inside bool

	Width int
	Height int

	shiftX int
	shiftY int
}

func NewSceneTrigger(width, height int) *SceneTrigger {
	var sceneTrigger *SceneTrigger = new(SceneTrigger)

	sceneTrigger.Width = width
	sceneTrigger.Height = height

	return sceneTrigger
}

func (this *SceneTrigger) Update(game *Game, dood *MultiSprite) {
	if (dood.X + 14 > float64(this.Width + this.shiftX * this.Width)) {
		if this.inside {
			return
		}
		
		this.shiftX++
		this.inside = true
	} else if (dood.X + 18 < float64(this.shiftX * this.Width)) {
		if this.inside {
			return
		}

		this.shiftX--
		this.inside = true
	} else if (dood.Y > float64(this.Height + this.shiftY * this.Height)) {
		if this.inside {
			return
		}

		this.shiftY++
		this.inside = true
	} else if (dood.Y + 32 < float64(this.shiftY * this.Height)) {
		if this.inside {
			return
		}

		this.shiftY--
		this.inside = true
	} else {
		this.inside = false
	}
	if (this.inside) {
		game.MakeSceneTransition(this.shiftX * this.Width, this.shiftY * this.Height)		
	}

}
