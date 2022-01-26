package game

import "fmt"

//献图：翻牌技：对方军营回复6点生命
type SkillXianTu struct {
	HeroSkill
}

func (ss *SkillXianTu) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (s *SkillXianTu) OnFaceUp(card *Card) {
	var xianTuValue int32 = 6
	g := card.GetPlayer().game

	for _, p := range g.players {
		if p.seatId != card.owner.seatId {
			//ar := NewActionRecoverCamp(g, p, card, xianTuValue)
			//ar.DoRecover()
			oldHP := p.hp
			p.hp += xianTuValue
			if p.hp > g.config.CampHP {
				p.hp = g.config.CampHP
			}
			SyncCampChangeHP(p, card, oldHP, p.hp)

			card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("张松 献图：翻牌技：对方军营回复6点生命（不超过上限）,实际回复%v", p.hp-oldHP))
			break
		}
	}

}
