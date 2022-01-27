package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/gameutil"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//2	曹昂	殿后	翻牌技：选择一名敌方武将，双方同时受到X点伤害。（X为生命值较低武将的生命值）

type SkillDianHou struct {
	HeroSkill

	targetCard *Card
}

func (ss *SkillDianHou) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Enemy

	//没有可选牌，直接发动
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}

	var minSelect int32 = 1
	var maxSelect int32 = 1

	if ret, ok := IsSelectTargetUnique(ss.card, selectCardType, maxSelect); ok {
		ss.targetCard = ret[0]
		return ret, nil
	}

	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌
	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCard = cards[0]
	return cards, nil
}

func (ss *SkillDianHou) OnFaceUp(card *Card) {
	if ss.targetCard == nil {
		return
	}

	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v发动%v,目标为%v", ss.card.GetOwnInfo(), ss.card.skillCfg.Name, ss.targetCard.GetOwnInfo()))
	damage := gameutil.Min(card.GetHP(), ss.targetCard.GetHP())

	g := card.owner.game
	g.PostActData(ss)
	ss.PostActStream(func() {
		NewActionDamageCard(g, card, card, nil, damage, card.GetSkillId()).DoDamage()
	})

	ss.PostActStream(func() {
		NewActionDamageCard(g, ss.targetCard, card, nil, damage, card.GetSkillId()).DoDamage()
	})
}
