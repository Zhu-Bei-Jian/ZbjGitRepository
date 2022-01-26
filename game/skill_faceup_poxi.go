package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//破袭：翻牌技，对一名武将造成8点伤害
type SkillPoXi struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillPoXi) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_other

	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}

	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌
	var minSelect int32 = 1
	var maxSelect int32 = 1

	if ret, ok := IsSelectTargetUnique(ss.card, selectCardType, maxSelect); ok {
		ss.targetCards = ret
		return ret, nil
	}

	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCards = cards
	return cards, nil
}

func (ss *SkillPoXi) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (ss *SkillPoXi) OnFaceUp(card *Card) {
	if ss.targetCards == nil {
		return
	}
	g := card.owner.game
	g.PostActData(ss)
	ss.PostActStream(func() {
		ad := NewActionDamageCard(g, ss.targetCards[0], card, nil, 8, ss.GetSkillId())
		ad.DoDamage()
	})
}
