package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//69	刘禅	享乐	翻牌技：使一名其他武将+5+5
type SkillFangQuan struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillFangQuan) CanUse() ([]*Card, error) {
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
func (s *SkillFangQuan) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}

	card.owner.game.PostActData(s)

	tCard := s.targetCards[0]
	card.owner.game.PostActData(s)
	s.PostActStream(func() {
		StartGetBuff(tCard, s.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
	})

}
