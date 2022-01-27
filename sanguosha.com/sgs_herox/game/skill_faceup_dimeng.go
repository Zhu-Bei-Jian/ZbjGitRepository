package game

import (
	"fmt"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

// 25	鲁肃	缔盟	翻牌技：指定一名其他武将获得此效果：该武将受到武将反击伤害-3。

type SkillDiMeng struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillDiMeng) CanUse() ([]*Card, error) {
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

	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCards = cards
	return cards, nil
}

func (ss *SkillDiMeng) OnFaceUp(card *Card) {
	if ss.targetCards == nil {
		return
	}
	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v发动%v,目标为%v", ss.card.GetOwnInfo(), ss.card.skillCfg.Name, ss.targetCards[0].GetOwnInfo()))
	StartGetBuff(ss.targetCards[0], ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
}
