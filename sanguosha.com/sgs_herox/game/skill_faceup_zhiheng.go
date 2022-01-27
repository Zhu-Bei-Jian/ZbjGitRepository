package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//37	孙权	制衡	翻牌技：使一名其他非重装正面武将牌翻面，且该武将当前回合不能翻牌
type SkillZhiHeng struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillZhiHeng) PreUseSkill() {

	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})

}

func (ss *SkillZhiHeng) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_NotHeavy

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

func (s *SkillZhiHeng) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}
	if s.targetCards[0].isBack {
		return
	}

	NewActionFaceDown(card.owner.game, s.targetCards[0], card).DoFaceDown()

}
