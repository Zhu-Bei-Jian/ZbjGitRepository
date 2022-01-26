package game

//效果：改攻击距离
type BuffCommonAttackDistance struct {
	distance int32
	HeroBuff
}

func (ss *BuffCommonAttackDistance) OnEnable(card *Card) {
	setAttackDistanceAndNotify(card, ss.distance)
}

func (ss *BuffCommonAttackDistance) OnDisable(card *Card) {
	setAttackDistanceAndNotify(card, 1)
}
