package game

import "sanguosha.com/sgs_herox/proto/cmsg"

type ActionFaceDown struct {
	ActionDataBase

	card *Card

	spellCard *Card
}

func NewActionFaceDown(g *GameBase, card *Card, spellCard *Card) *ActionFaceDown {
	ac := &ActionFaceDown{}
	ac.game = g
	ac.card = card
	ac.spellCard = spellCard
	return ac
}

func (ac *ActionFaceDown) DoFaceDown() {

	ac.game.PostActData(ac)

	ac.PostActStream(func() {
		ac.card.isBack = true     //翻回背面
		ac.card.canTurnUp = false //本回合不能翻至正面
		if ac.card.skillId != -1 {
			ac.card.skillId = 0 //自身技能清除   , 但其他武将对其施加的技能和buff 不清除
		}

		ac.card.GetPlayer().GetGame().Send2All(&cmsg.SSyncCard{
			SeatId:    ac.card.GetPlayer().GetSeatID(),
			Card:      ac.card.ToDef(ac.card.GetPlayer().GetSeatID()),
			SpellCard: ac.spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_FaceDown,
			},
			}})
	})

	ac.PostActStream(func() {
		if ac.card.skillId == -1 {
			return
		}
		ac.game.GetTriggerMgr().Trigger(TriggerType_LoseSkill, ac, ac.card)
	})

}
