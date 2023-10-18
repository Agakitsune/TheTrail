package engine

type SceneTrigger struct {
	inside bool

	Width int
	Height int

	ShiftX int
	ShiftY int
}

func NewSceneTrigger(width, height int) *SceneTrigger {
	var sceneTrigger *SceneTrigger = new(SceneTrigger)

	sceneTrigger.Width = width
	sceneTrigger.Height = height

	return sceneTrigger
}

func (this *SceneTrigger) Update(game *Game, dood *MultiSprite) {
	if (dood.X + 14 > float64(this.Width + this.ShiftX * this.Width)) {
		if this.inside {
			return
		}
		
		this.ShiftX++
		this.inside = true
	} else if (dood.X + 18 < float64(this.ShiftX * this.Width)) {
		if this.inside {
			return
		}

		this.ShiftX--
		this.inside = true
	} else if (dood.Y > float64(this.Height + this.ShiftY * this.Height)) {
		if this.inside {
			return
		}

		this.ShiftY++
		this.inside = true
	} else if (dood.Y + 32 < float64(this.ShiftY * this.Height)) {
		if this.inside {
			return
		}

		this.ShiftY--
		this.inside = true
	} else {
		this.inside = false
	}
	if (this.inside) {
		game.MakeSceneTransition(this.ShiftX * this.Width, this.ShiftY * this.Height)		
	}

}
