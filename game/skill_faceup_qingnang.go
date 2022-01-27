package game

import (
	"fmt"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

type SkillQingNang struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillQingNang) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Any

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

//17	华佗	青囊	翻牌技：指定一名武将恢复10点生命值
func (ss *SkillQingNang) OnFaceUp(card *Card) {
	if ss.targetCards == nil {
		return
	}
	t := ss.targetCards[0]
	oldHp := t.GetHP()
	t.AddHP(ss.GetValue(1))
	SyncChangeHP(t, oldHp, t.GetHP(), card, ss.GetSkillId())
	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 青囊 翻牌技：指定一名武将(%v)恢复%v点生命值", card.GetOwnInfo(), t.GetOwnInfo(), ss.GetValue(1)))
}
