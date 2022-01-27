package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

//凶镬：翻牌技：对敌置牌区排所有武将造成3点伤害
type SkillXiongHuo struct {
	HeroSkill
}

func (ss *SkillXiongHuo) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (ss *SkillXiongHuo) OnFaceUp(card *Card) {
	var enemyRow int32 //敌方后排
	if card.owner.seatId == 1 {
		enemyRow = 0
	} else {
		enemyRow = 2
	}
	card.owner.game.PostActData(ss)
	ss.PostActStream(func() {
		card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动 凶镬：翻牌技：对敌置牌区排所有武将造成%v点伤害", card.GetOwnInfo(), ss.GetValue(1)))
		logrus.Info("-------开始扫描敌方置牌区--------")
	})
	ss.PostActStream(func() {
		for i := 0; i < 3; i++ {
			ti := i
			if cd, ok := card.owner.game.board.cells[enemyRow][ti].GetCard(); ok {
				ss.PostActStream(func() {
					card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v受到凶镬的3点伤害", cd.GetOwnInfo()))
					NewActionDamageCard(card.owner.game, cd, card, nil, ss.GetValue(1), ss.GetSkillId()).DoDamage()
				})
			}
		}
	})

}
