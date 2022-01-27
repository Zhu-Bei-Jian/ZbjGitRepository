package game

//谦逊：翻牌技，攻击力减少2。
type SkillQianXun struct {
	HeroSkill
}

func (ss *SkillQianXun) OnFaceUp(card *Card) {
	card.attack -= 2
	SyncChangeAttack(card, card.attack+2, card.attack, card)
}
