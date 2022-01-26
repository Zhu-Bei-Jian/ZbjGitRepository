package game

import (
	"fmt"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//激将：翻牌技，令一名武将的攻击力+2，若自身在交战区，可令两名武将攻击力+2。
type SkillJiJiang struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillJiJiang) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Any
	selectCount := int32(1)
	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}
	if IsInWarZone(ss.card) {
		selectCount++
	}
	if ret, ok := IsSelectTargetUnique(ss.card, selectCardType, selectCount); ok {
		ss.targetCards = ret
		return ret, nil
	}
	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌
	var minSelect int32 = 1
	var maxSelect int32 = 2

	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCards = cards
	return cards, nil
}

func (ss *SkillJiJiang) OnFaceUp(card *Card) {
	if ss.targetCards == nil {
		return
	}
	//Log(fmt.Sprintf("%v 发动激将：令一名武将的攻击力+2，若自身在交战区，可令两名武将攻击力+2。", card.GetOwnInfo()))
	g := card.owner.game
	g.PostActData(ss)

	for _, cd := range ss.targetCards {
		t := cd
		ss.PostActStream(func() {
			StartGetBuff(t, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
		})
		ss.PostActStream(func() {
			t.attack += 2
			card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 受到激将加成。攻击力+2", t.GetOwnInfo()))
			SyncChangeAttack(t, t.attack-2, t.attack, t)
		})

	}
}
