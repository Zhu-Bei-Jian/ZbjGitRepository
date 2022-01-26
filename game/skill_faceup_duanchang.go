package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//沉默一个武将，若在交战区，沉默2个角色
type SkillDuanChang struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillDuanChang) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Any

	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}
	var minSelect int32 = 1
	var maxSelect int32 = 1
	if IsInWarZone(ss.card) {
		maxSelect = 2
	}

	if ret, ok := IsSelectTargetUnique(ss.card, selectCardType, maxSelect); ok {
		ss.targetCards = ret
		return ret, nil
	}
	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌

	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCards = cards
	return cards, nil
}
func (s *SkillDuanChang) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}

	card.owner.game.PostActData(s)
	for _, c := range s.targetCards {
		t := c
		s.PostActStream(func() {
			StartSilentCard(t, card)
		})
	}

}

//沉默一张卡牌
func StartSilentCard(card *Card, spellCard *Card) {
	ad := NewActionSilence(card.owner.game, card, spellCard)
	ad.DoSilence()
}
