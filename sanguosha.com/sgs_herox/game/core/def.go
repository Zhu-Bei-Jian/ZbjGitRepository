package core

type Worker interface {
	Post(func())
}

type Player interface {
}
