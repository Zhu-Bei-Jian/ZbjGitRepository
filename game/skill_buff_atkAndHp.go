package game

//修改 血量 （会直接修改上限） 的buff
//修改 攻击力
type BuffAtkAndHp struct {
	attack int32
	hp     int32
	HeroBuff
}

func (ss *BuffAtkAndHp) SetAttack(attack int32) {
	ss.attack = attack
}

func (ss *BuffAtkAndHp) SetHP(hp int32) {
	ss.hp = hp
}

func (ss *BuffAtkAndHp) OnEnable(card *Card) {

	if ss.hp > 0 {
		oldHP := card.GetHP()
		card.AddHpMax(ss.hp)
		card.AddHP(ss.hp)
		SyncChangeHP(card, oldHP, card.GetHP(), card, card.GetSkillId())
	} else if ss.hp < 0 {
		oldHP := card.GetHP()
		card.SubHP(ss.hp, false)
		SyncChangeHP(card, oldHP, card.GetHP(), card, card.GetSkillId())
	}

	ModifyHPAttack(card, 0, ss.attack, nil, ss.GetBuffId())

}

func (ss *BuffAtkAndHp) OnDisable(card *Card) {
	//暂时不需要还原血量
	//........

	//还原攻击力
	ModifyHPAttack(card, 0, -ss.attack, nil, ss.GetBuffId())

}
