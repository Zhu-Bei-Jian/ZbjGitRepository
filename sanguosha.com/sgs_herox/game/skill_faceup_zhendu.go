package game

import gamedef "sanguosha.com/sgs_herox/proto/def"

//鸩毒：翻牌技：对己方一名其他明牌武将造成3点伤害
type SkillZhenDu struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillZhenDu) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_MyOwnFaceUp

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

func (ss *SkillZhenDu) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (s *SkillZhenDu) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}
	g := card.owner.game
	g.PostActData(s)
	s.PostActStream(func() {
		ad := NewActionDamageCard(g, s.targetCards[0], card, nil, s.GetValue(1), s.GetSkillId())
		ad.DoDamage()
	})
}
