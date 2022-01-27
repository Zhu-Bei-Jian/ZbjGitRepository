package game

import (
	"fmt"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//反间：翻牌技，对一名其它友方武将造成3点伤害。
type SkillFanJian struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillFanJian) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_OtherMyOwn

	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		ss.card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v %v。检测场上没有可选目标，不发动此技能，但仍翻至正面.", ss.card.heroCfg.Name, ss.skillCfg.Name))
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

func (s *SkillFanJian) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}

	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v发动%v,目标为%v", s.card.GetOwnInfo(), s.card.skillCfg.Name, s.targetCards[0].GetOwnInfo()))
	p := card.GetPlayer()
	g := p.game
	target := s.targetCards[0]

	actDamage := NewActionDamageCard(g, target, card, nil, s.GetValue(1), s.GetSkillId())
	actDamage.DoDamage()

}
