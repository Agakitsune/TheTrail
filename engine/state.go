package engine

type State interface {
	Load(*Game)
	Update()
	Draw()
}
