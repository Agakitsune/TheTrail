package engine

import (
	"image"
)

type Animation struct {
	Frames    []int
	FrameSize Vector2
	Row       int

	LoopOn int

	Selection int
	Speed     int

	Count int
}

type Animator struct {
	Animations map[string]*Animation

	Current string
}

func (this *Animator) SetFrameSize(x float64, y float64) {
	for _, anim := range this.Animations {
		anim.SetFrameSize(x, y)
	}
}

func (this *Animation) SetFrameSize(x float64, y float64) {
	this.FrameSize = Vector2{X: x, Y: y}
}

func (this *Animation) GetFrame(frame int) int {
	return this.Frames[frame]
}

func (this *Animation) Update(sprite *MultiSprite) {
	this.Count++

	if this.Count%this.Speed != 0 {
		return
	}
	this.Selection++

	if this.Selection >= len(this.Frames) {
		this.Selection = this.LoopOn
	}

	frame := this.GetFrame(this.Selection)
	sprite.Rect = image.Rect(frame*int(this.FrameSize.X),
		this.Row*int(this.FrameSize.Y),
		frame*int(this.FrameSize.X)+int(this.FrameSize.X),
		this.Row*int(this.FrameSize.Y)+int(this.FrameSize.Y))
}

func (this *Animator) SetAnimation(name string) {
	if this.Current == name {
		return
	}
	var anim = this.Animations[this.Current]
	anim.Selection = 0
	this.Current = name
	anim = this.Animations[this.Current]
	anim.Selection = 0
}

func (this *Animator) OnLastFrame() bool {
	return this.Animations[this.Current].Selection == len(this.Animations[this.Current].Frames)-1
}

func (this *Animator) Update(sprite *MultiSprite) {
	var anim = this.Animations[this.Current]
	anim.Update(sprite)
}
