package game

import (
	"fmt"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//断粮：翻牌技，对一名武将造成4点伤害，并与其互换位置
type SkillDuanLiang struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillDuanLiang) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Any

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

func (s *SkillDuanLiang) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}

	p := card.GetPlayer()
	g := p.game
	target := s.targetCards[0]
	targetCell := s.targetCards[0].cell
	g.PostActData(s)
	s.PostActStream(func() {
		actDamage := NewActionDamageCard(g, target, card, nil, 4, s.GetSkillId())
		actDamage.DoDamage()
	})

	s.PostActStream(func() {
		if !target.IsDead() {
			card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("徐晃发动技能断粮，对%v造成4点伤害,%v死亡", target.heroCfg.Name, target.heroCfg.Name))
			StartExchangePos(card, target)
		} else {
			card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("徐晃发动技能断粮，对%v造成4点伤害", target.heroCfg.Name))
			StartMoveToCell(card, targetCell)
		}

	})

}
