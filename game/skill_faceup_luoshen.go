package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//59	甄姬	洛神	翻牌技：指定一个友方武将本回合攻击力+3
type SkillLuoShen struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillLuoShen) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_MyOwn
	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}
	if ret, ok := IsSelectTargetUnique(ss.card, selectCardType, 1); ok {
		ss.targetCards = ret
		return ret, nil
	}
	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌
	var minSelect int32 = 1
	var maxSelect int32 = 1

	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCards = cards
	return cards, nil
}

func (ss *SkillLuoShen) OnFaceUp(card *Card) {
	if ss.targetCards == nil {
		return
	}
	StartGetBuff(ss.targetCards[0], ss.GetBuffId0(), gameconf.ExpireTyp_ETRound, 1, card)
}
