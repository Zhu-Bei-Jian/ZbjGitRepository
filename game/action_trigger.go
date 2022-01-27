package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type TriggerType int32

const (
	TriggerOpp_None        TriggerType = iota
	TriggerType_PhaseBegin             //回合开始时
	TriggerType_DrawStart              //抓牌开始时
	TriggerType_PhaseMain              //操作开始时
	TriggerType_PhaseEnd               //结束阶段开始时

	TriggerType_BeforeActAttack //攻击行动前

	TriggerType_AttackCard     //攻击单张牌
	TriggerType_BeAttackCard   //被攻击时
	TriggerType_BeRetAttack    //受到反击
	TriggerType_MakeDamage     //造成伤害时
	TriggerType_OnDamage       //受到伤害时
	TriggerType_AfterBeDamaged //受到伤害后
	TriggerType_AfterAttack    //攻击后

	TriggerType_AttackCamp               //攻击军营
	TriggerType_BeRetAttackByCamp        //军营反击伤害
	TriggerType_MakeDamageCamp           //对军营造成伤害时
	TriggerType_CampMakeDamageWillEffect //军营对卡造成伤害前
	TriggerType_CampMakeDamage           //军营对卡造成伤害
	TriggerType_AfterAttackCamp          //攻击军营后

	TriggerType_MoveCard       //移动卡牌时
	TriggerType_AfterPosChange //卡牌位置有变化时
	TriggerType_AfterMoveCard  //移动卡牌后
	TriggerType_PlaceCard      //放置卡牌时

	TriggerType_GetSkill  //获得技能时
	TriggerType_LoseSkill //失去技能时

	TriggerType_CheckHeavyAP
	TriggerType_CheckAP
	TriggerType_CheckDistance
	TriggerType_CheckAttackCntInTurn
	TriggerType_CheckAttackCamp

	TriggerType_CheckAP_PlaceCard

	TriggerType_Die //卡牌死亡后
)

type TriggerHandler struct {
	name         string
	triggerTypes []TriggerType
	handle       func(g *GameBase, ad core.IActionData, params ...interface{})
}

type TriggerManager struct {
	game     *GameBase
	triggers map[TriggerType][]TriggerHandler
}

func newTriggerManager(g *GameBase) *TriggerManager {
	return &TriggerManager{triggers: make(map[TriggerType][]TriggerHandler), game: g}
}

func (tm *TriggerManager) Register(handler TriggerHandler) {
	for _, triggerType := range handler.triggerTypes {
		if _, ok := tm.triggers[triggerType]; ok {
			tm.triggers[triggerType] = append(tm.triggers[triggerType], handler)
		} else {
			tm.triggers[triggerType] = []TriggerHandler{handler}
		}
	}
}

func (tm *TriggerManager) init() {
	skillCfgs := tm.game.config.Skill.All(gameconf.SkillTyp_STPassive)
	for _, skillCfg := range skillCfgs {
		skill, exist := NewSkill(skillCfg)
		if !exist {
			continue
		}
		handlers := skill.TriggerHandler()
		for _, handler := range handlers {
			tm.Register(handler)
		}
	}

	buffCfgs := tm.game.config.Buff.All()
	for _, buffCfg := range buffCfgs {
		tBuff, exist := newBuff(buffCfg)
		if !exist {
			continue
		}
		handlers := tBuff.TriggerHandler()
		for _, handler := range handlers {
			tm.Register(handler)
		}
	}

	//注册 buff_common
	t, exist := newBuff(&gameconf.BuffConfDefine{})
	if !exist {
		return
	}
	handlers := t.TriggerHandler()
	for _, handler := range handlers {
		tm.Register(handler)
	}
}

type ActionTrigger struct {
	ActionDataBase
}

func NewActionTrigger(game *GameBase) *ActionTrigger {
	ac := &ActionTrigger{}
	ac.game = game
	return ac
}

func (tm *TriggerManager) Trigger(triggerType TriggerType, core core.IActionData, params ...interface{}) {
	ac := NewActionTrigger(tm.game)
	tm.game.PostActData(ac)

	ac.PostActStream(func() {
		if v, ok := tm.triggers[triggerType]; ok {
			for _, tr := range v {
				tr := tr
				ac.PostActStream(func() {
					tr.handle(ac.game, core, params...)
				})
			}
		}
	})

	return
}
