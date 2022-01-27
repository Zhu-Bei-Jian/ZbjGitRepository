package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//节命：翻牌技，令一名其他武将免疫下一次受到的伤害，若自身在敌方置牌区，额外选择一名武将。
type SkillJieMing struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillJieMing) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_other

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
	if IsInEnemyZone(ss.card) {
		maxSelect = 2
	}
	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCards = cards
	return cards, nil
}

func (ss *SkillJieMing) OnFaceUp(card *Card) {
	if ss.targetCards == nil {
		return
	}
	g := card.owner.game
	cfg := g.GetBuffCfg(ss.GetBuffId0())
	for _, v := range ss.targetCards {
		StartGetBuff(v, ss.GetBuffId0(), gameconf.ExpireTyp_ETTimes, cfg.GetExpireV(), card)
	}
}
