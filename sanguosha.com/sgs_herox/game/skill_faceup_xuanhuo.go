package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//10	法正	眩惑	翻牌技，令一名武将攻击力和生命值互换
type SkillXuanHuo struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillXuanHuo) CanUse() ([]*Card, error) {
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
func (ss *SkillXuanHuo) OnFaceUp(card *Card) {
	if ss.targetCards == nil {
		return
	}
	//眩惑 会影响血量上限 ：交换后的血量值 ，同时作为卡牌的血量上限
	card.owner.game.PostActData(ss)
	tCard := ss.targetCards[0]
	oldHp := tCard.GetHP()
	oldAttack := tCard.attack

	//区分交换后 血量上限是增加还是减少，分开处理
	if oldHp <= oldAttack { //增加
		tCard.AddHpMax(oldAttack - oldHp)
		tCard.AddHP(oldAttack - oldHp)
	} else { //减少
		tCard.SubHpMax(oldHp - oldAttack)
	}

	tCard.attack = oldHp

	ss.PostActStream(func() {
		StartGetBuff(tCard, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
	})
	ss.PostActStream(func() {
		SyncChangeAttack(tCard, tCard.attack, oldAttack, card)
		SyncChangeHP(tCard, tCard.GetHP(), oldHp, card, ss.GetSkillId())
	})

}
