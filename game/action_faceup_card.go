package game

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type ActionFaceUpCard struct {
	ActionDataBase

	card      *Card
	isPassive bool
	spellCard *Card

	actSelect *gamedef.ActSelectParam
}

func NewActionFaceUpCard(g *GameBase, card *Card, isPassive bool, spellCard *Card, actSelect *gamedef.ActSelectParam) *ActionFaceUpCard {
	ac := &ActionFaceUpCard{}
	ac.game = g
	ac.card = card
	ac.isPassive = isPassive
	ac.actSelect = actSelect
	return ac
}
func (ac *ActionFaceUpCard) Error() string {
	return "ActionFaceUpCard Error"
}
func (ac *ActionFaceUpCard) DoFaceUp(cb func(err error)) {
	ac.game.PostActData(ac)

	card := ac.card

	//为符合客户端显示逻辑 ，调整同步顺序 为  同步faceup，同步buff，同步hp和attack

	ac.PostActStream(func() {

		//skillId = 0  代表背面原始状态 ，无技能 但没被沉默
		//skillId = -1 代表遭到了沉默，
		//被动翻牌不触发翻牌技能
		//翻牌技能触发
		flag := false
		if !ac.isPassive && card.skillCfg.SkillType == gameconf.SkillTyp_STFaceUp && card.GetSkillId() == 0 && card.isBack { //card没有被沉默 且处于背面，才能发动技能
			flag = true
			card.skillId = card.skillCfg.SkillID //可以发动技能，对 skillId 进行赋值(Id>0)
			skill, ok := NewSkill(card.skillCfg)
			if !ok {
				logrus.Errorf("new skill %d fail", card.skillCfg.SkillID)
				return
			}
			skill.SetActSelect(ac.actSelect)
			skill.SetCard(card)
			targets, err := skill.CanUse()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"skillId":   card.skillId,
					"skillName": card.skillCfg.Name,
				}).WithError(err).Error("skill.CanUse")
				ac.Stop()
				if cb != nil {
					cb(err)
				}
				card.skillId = 0

				return
			}

			card.isBack = false
			if card.skillId != -1 {
				card.skillId = card.skillCfg.SkillID
			}
			g := card.owner.game
			//var targetsToDef []*gamedef.Card
			//for _, v := range targets {
			//	targetsToDef = append(targetsToDef, v.ToDef(-1))
			//}
			skill.SetTargets(targets)

			ac.PostActStream(func() {
				//g.Send2All(&cmsg.SSyncUseSkill{
				//	SkillId:    card.skillCfg.SkillID,
				//	Card:       card.ToDef(-1),
				//	TargetCard: targetsToDef,
				//})
				g.Send2All(&cmsg.SSyncCard{
					SeatId:    card.GetPlayer().GetSeatID(),
					Card:      card.ToDef(-1),
					SpellCard: ac.spellCard.ID(),
					Changes: []*cmsg.SSyncCard_Change{{
						ChangeType: cmsg.SSyncCard_FaceUp,
					},
					}})
			})
			if targets != nil {
				ac.PostActStream(func() {
					skill.PreUseSkill()
				})
			}

			ac.PostActStream(func() {
				skill.OnFaceUp(card)
			})

		}
		if !flag {
			card.isBack = false
			if card.skillId != -1 {
				card.skillId = card.skillCfg.SkillID
			}
			card.GetPlayer().GetGame().Send2All(&cmsg.SSyncCard{
				SeatId:    card.GetPlayer().GetSeatID(),
				Card:      card.ToDef(-1),
				SpellCard: ac.spellCard.ID(),
				Changes: []*cmsg.SSyncCard_Change{{
					ChangeType: cmsg.SSyncCard_FaceUp,
				},
				}})
		}

	})

	ac.PostActStream(func() {
		//被动技才会有获得技能触发
		if card.skillCfg.SkillType == gameconf.SkillTyp_STPassive && card.skillId > 0 {
			ac.game.GetTriggerMgr().Trigger(TriggerType_GetSkill, ac, card)
		}
	})
	ac.PostActStream(func() {
		if cb != nil {
			cb(nil)
		}
	})
	ac.PostActStream(func() {
		ac.game.TryEndGame()
	})
}
