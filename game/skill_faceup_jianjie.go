package game

import (
	"fmt"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//荐杰：翻牌技,令一名武将身材变为10/4  (attack/hp)
type SkillJianJie struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillJianJie) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Any

	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		ss.card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 翻牌技 %v。检测场上没有可选目标，不发动此技能，但仍翻至正面.", ss.card.heroCfg.Name, ss.skillCfg.Name))
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

func (s *SkillJianJie) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}
	card.owner.game.PostActData(s)
	s.PostActStream(func() {
		StartGetBuff(s.targetCards[0], s.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
	})
	s.PostActStream(func() {
		ModifyHPAttack(s.targetCards[0], -s.targetCards[0].GetHP()+4, -s.targetCards[0].attack+10, card, s.GetSkillId())

	})

}
