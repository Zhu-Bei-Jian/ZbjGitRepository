package game

import gamedef "sanguosha.com/sgs_herox/proto/def"

type CheckActType = int32

const (
	CHECK_NONE  CheckActType = 101
	MOVE_CARD   CheckActType = 102
	PLACE_CARD  CheckActType = 103
	ATTACK_CARD CheckActType = 104
	ATTACK_CAMP CheckActType = 105
)

//检查能否行动
type ActionCheck struct {
	ActionDataBase

	actType    gamedef.ActType
	card       *Card
	targetCard *Card

	ap              map[int32]int32 // ap[MOVE_CARD] 本次移动卡牌所要耗费的行动点 ，依此类推
	maxDistance     int32
	maxAttackCount  int32
	canAttackCamp   bool
	maxCampDistance int32 //攻击军营距离

	effectSkills []Skill
}

func NewActionCheck(g *GameBase, card *Card, actType gamedef.ActType, targetCard *Card) *ActionCheck {
	ac := &ActionCheck{}
	ac.game = g
	ac.card = card
	ac.targetCard = targetCard
	ac.actType = actType
	ac.ap = make(map[int32]int32)
	ac.ap[MOVE_CARD] = 1
	ac.ap[PLACE_CARD] = 1
	ac.ap[ATTACK_CARD] = 1
	ac.ap[ATTACK_CAMP] = 1
	ac.canAttackCamp = true
	ac.maxCampDistance = 1
	ac.maxAttackCount = 1
	return ac
}

func (ac *ActionCheck) DoCheck(callback func()) {
	ac.game.PostActData(ac)

	//ac.PostActStream(func() {
	//	ac.game.GetTriggerMgr().Trigger(TriggerType_CheckAP, ac)
	//})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_CheckDistance, ac)
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_CheckAttackCntInTurn, ac)
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_CheckAttackCamp, ac)
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_CheckHeavyAP, ac)
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_CheckAP, ac)
	})
 
	ac.PostActStream(func() {
		callback()
		return
	})
}

func (ac *ActionCheck) AddEffectSkill(skill Skill) {
	ac.effectSkills = append(ac.effectSkills, skill)
}
