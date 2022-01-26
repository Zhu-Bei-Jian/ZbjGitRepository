package game

//效果：改攻击次数的上限
type BuffCommonMaxAttackCount struct {
	maxCount int32
	HeroBuff
}

func (ss *BuffCommonMaxAttackCount) OnEnable(card *Card) { //国色，令一名武将 攻击次数 无限制
	setAttackCountAndNotify(card, ss.maxCount)
}

func (ss *BuffCommonMaxAttackCount) OnDisable(card *Card) { //一般
	setAttackCountAndNotify(card, 1) //默认上限为1
}
