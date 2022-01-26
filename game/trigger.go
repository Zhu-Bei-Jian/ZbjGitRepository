package game

type Register interface {
	RegisterTrigger()
}

func RegisterTrigger(register Register) {
	register.RegisterTrigger()
}
