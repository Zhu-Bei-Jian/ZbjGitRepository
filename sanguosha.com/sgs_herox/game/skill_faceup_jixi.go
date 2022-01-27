package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//急袭：翻牌技，交换任意两名武将的位置。
type SkillJiXi struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillJiXi) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Any

	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}

	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌
	var minSelect int32 = 2
	var maxSelect int32 = 2

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

func (s *SkillJiXi) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}
	//Log(fmt.Sprintf("%v 发动急袭 ，令%v 与%v 交换", card.GetOwnInfo(), s.targetCards[0].GetOwnInfo(), s.targetCards[1].GetOwnInfo()))
	StartExchangePos(s.targetCards[0], s.targetCards[1])
}
