package game

import gamedef "sanguosha.com/sgs_herox/proto/def"

// 离间	翻牌技：选择1名其他友方武将攻击1名敌方武将直至一方阵亡。

type SkillLiJian struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillLiJian) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_OneOtherMyOwnAndOneEnemy

	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}

	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌
	var minSelect int32 = 2
	var maxSelect int32 = 2

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
func (s *SkillLiJian) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}
	g := card.owner.game
	g.PostActData(s)
	myCard := s.targetCards[0]
	enemy := s.targetCards[1]
	//离间 需要 双方先翻面（被动） //避免出现 0伤 不翻面
	s.PostActStream(func() {
		if myCard.isBack {
			NewActionFaceUpCard(g, myCard, true, card, nil).DoFaceUp(nil)
		}
	})
	s.PostActStream(func() {
		if enemy.isBack {
			NewActionFaceUpCard(g, enemy, true, card, nil).DoFaceUp(nil)
		}
	})
	var f func(int32)
	f = func(cnt int32) {
		s.PostActStream(func() {
			if enemy.IsDead() || myCard.IsDead() || cnt > 100 {
				s.Stop()
			}
		})

		s.PostActStream(func() {
			NewActionActAttack(g, myCard, []*Card{enemy}, card, nil).DoAttack()
		})

		s.PostActStream(func() {
			g.StartWaitingNoneFloat(1.4, nil)
			f(cnt + 1)
		})

	}
	f(1)
	//加入决斗次数计数，超过一定上限，视为双方都是零伤害的死循环，强制退出

}
