package engine

import(
	"image"
)

type Animation struct {
	frames []int
	row int

	loopOn int

	selection int
	speed int

	count int
}

type Animator struct {
	animations map[string]*Animation

	current string
}

func (this *Animation) GetFrame(frame int) int {
	return this.frames[frame]
}

func (this *Animation) Update(sprite *MultiSprite) {
	this.count++

	if this.count % this.speed != 0 {
		return
	}
	this.selection++

	if this.selection >= len(this.frames) {
		this.selection = this.loopOn
	}

	frame := this.GetFrame(this.selection)
	sprite.rect = image.Rect(frame * 32, this.row * 32, frame * 32 + 32, this.row * 32 + 32)
}

func (this *Animator) SetAnimation(name string) {
	if this.current == name {
		return
	}
	var anim = this.animations[this.current]
	anim.selection = 0
	this.current = name
	anim = this.animations[this.current]
	anim.selection = 0
}

func (this *Animator) Update(sprite *MultiSprite) {
	var anim = this.animations[this.current]
	anim.Update(sprite)
}
