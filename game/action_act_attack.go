package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

type ActionActAttack struct {
	ActionDataBase

	srcCard      *Card
	targetCards  []*Card
	spellCard    *Card   //引发这次攻击的卡牌  //一般情况下  spellCard=srcCard  ,貂蝉发动离间时,spellCard =貂蝉
	effectSkills []Skill //生效的技能

	extraRetDamage int32 //额外反击伤害
	data           []int32
}

func NewActionActAttack(g *GameBase, srcCard *Card, targetCards []*Card, spellCard *Card, effectSkills []Skill) *ActionActAttack {
	ac := &ActionActAttack{}
	ac.game = g
	ac.srcCard = srcCard
	ac.targetCards = targetCards
	ac.spellCard = spellCard
	ac.effectSkills = effectSkills
	return ac
}

func (ac *ActionActAttack) DoAttack() {

	ac.game.PostActData(ac)

	for _, v := range ac.effectSkills {
		ac.PostActStream(func() {
			v.PreUseSkill()
		})
	}

	ac.PostActStream(func() {
		if ac.srcCard.ID() == ac.spellCard.ID() { //检测 是否为主动发起的攻击 （or 离间引起
			ac.srcCard.attackCntInTurn++
			ac.srcCard.owner.attackCount++
		}

		opCard := ac.srcCard
		targetPos := make([]*gamedef.Position, 0)
		for _, v := range ac.targetCards {
			targetPos = append(targetPos, v.cell.Position.ToDef())
		}

		if len(targetPos) == 1 && calDistanceByPosition(ac.srcCard.cell.Position, toPosition(targetPos[0])) > 1 && ac.srcCard.skillId == 20 && ac.spellCard.heroCfg.Name != "貂蝉" {
			return
		}

		var spellCard int32
		if ac.spellCard != nil {
			spellCard = ac.spellCard.id
		}
		ac.game.Send2All(&cmsg.SSyncAct{
			ActType:  gamedef.ActType_AttackCard,
			OpSeatId: opCard.GetPlayer().seatId,
			OpCard:   opCard.id,
			//OpPos:     opCard.cell.Position.ToDef(),
			TargetPos: targetPos,
			//Card:      opCard.ToDef(-1),
			SpellCard: spellCard,
		})
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_BeforeActAttack, ac)
	})
	var targetCell []*Cell
	ac.PostActStream(func() {
		for _, v := range ac.targetCards {
			targetCell = append(targetCell, v.cell)
		}
	})

	ac.PostActStream(func() {
		for _, v := range ac.targetCards {
			card := v
			ac.PostActStream(func() {
				if ac.srcCard.GetHP() <= 0 {
					ac.Stop()
					return
				}
				ad := NewActionAttackCard(ac.game, ac.srcCard, card)
				ad.DoAttack()
			})
		}
	})

	ac.PostActStream(func() {
		deadId := -1
		for i := len(ac.targetCards) - 1; i >= 0; i-- {
			if ac.targetCards[i].IsDead() {
				deadId = i
				break
			}
		}
		if deadId == -1 {
			return
		}
		if ac.srcCard.IsDead() {
			return
		}
		if calDistanceByPosition(ac.srcCard.cell.Position, targetCell[deadId].Position) > 1 && ac.srcCard.skillId == 20 && ac.spellCard.heroCfg.Name != "貂蝉" {
			return
		}

		ac.srcCard.owner.Log(fmt.Sprintf("%v位置替换至%v", ac.srcCard.GetOwnInfo(), targetCell[deadId].Position))
		fromPos := ac.srcCard.cell.Position
		card := ac.srcCard.cell.RemoveCard()
		ac.game.board.cells[targetCell[deadId].Row][targetCell[deadId].Col].SetCard(card)
		SyncChangePos(&fromPos, ac.game, card)

		ac.PostActStream(func() {
			ac.game.GetTriggerMgr().Trigger(TriggerType_AfterPosChange, ac)
		})
	})

}
