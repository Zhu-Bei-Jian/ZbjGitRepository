package game

import "fmt"

// 解烦：翻牌技，令己方军营回复8点生命或对敌方军营造成8点伤害

type SkillJieFan struct {
	HeroSkill
	targetSeat int32
}

func (ss *SkillJieFan) CanUse() ([]*Card, error) {
	ss.targetSeat = ss.actSelectParam.Seat
	return nil, nil
}

func (ss *SkillJieFan) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (s *SkillJieFan) OnFaceUp(card *Card) {

	selectCampId := s.targetSeat
	g := card.GetPlayer().game
	mySeatId := card.GetPlayer().seatId

	if selectCampId == mySeatId {
		oldhp := g.players[mySeatId].GetHP()
		if oldhp == g.config.CampHP {
			return
		}
		g.players[mySeatId].hp += 8
		if g.players[mySeatId].hp > g.config.CampHP {
			g.players[mySeatId].hp = g.config.CampHP
		}
		SyncCampChangeHP(g.players[mySeatId], card, oldhp, g.players[mySeatId].GetHP()-oldhp)
		card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动解烦。令己方军营回复8点生命 ", card.GetOwnInfo()))
	} else {
		//必须走伤害流程ActionDamageCamp，需要与其他技能交互（军营免疫伤害 等）
		NewActionDamageCamp(g, g.players[selectCampId], card, 8).DoDamage()
		card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动解烦。对敌方军营造成8点伤害", card.GetOwnInfo()))
	}

}
