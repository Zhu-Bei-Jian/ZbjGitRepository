package game

//效果：可翻正面否
type BuffCommonFaceUp struct {
	canFaceUp bool
	HeroBuff
}

func (ss *BuffCommonFaceUp) OnEnable(card *Card) {
	setCanFaceUpAndNotify(card, ss.canFaceUp)
}

func (ss *BuffCommonFaceUp) OnDisable(card *Card) {
	setCanFaceUpAndNotify(card, !ss.canFaceUp)
}
