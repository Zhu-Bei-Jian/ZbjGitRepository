package game

import (
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//检查能否行动
type ActionCheckPlaceCard struct {
	ActionDataBase

	heroCfg  *gameconf.HeroDefine
	skillCfg *gameconf.SkillConfDefine
	opPlayer *Player

	ap    int32           //此次放置卡牌 所需消耗的 行动点
	buffs []HeavyHaloBuff //如果此次放置 最终成功，那么被放置的武将 将直接获得 buffs   （用于解决重装武将的光环Buff相关问题 ：一旦上场，就需要开始向光环受益者施加buff）。 区别于翻面施加buff的时机点。重装不存在翻面
}

type EffectType int32

const (
	EffectType_None = iota
	EffectType_Card
	EffectType_camp
)

type HeavyHaloBuff struct { //重装武将的光环buff
	buff
	effectType   EffectType //光环受益对象 是卡牌 还是 军营
	targetCards  []*Card    //光环受益者 卡牌
	targetSeatID int32      //光环受益者 军营
}

func NewActionCheckPlaceCard(g *GameBase, opPlayer *Player, heroCfg *gameconf.HeroDefine, skillCfg *gameconf.SkillConfDefine) *ActionCheckPlaceCard {
	ac := &ActionCheckPlaceCard{}
	ac.game = g
	ac.opPlayer = opPlayer
	ac.heroCfg = heroCfg
	ac.skillCfg = skillCfg
	ac.ap = g.config.GetCommonCost() //非重装武将的一次行动花费
	return ac
}

func (ac *ActionCheckPlaceCard) DoCheck(callback func()) {
	ac.game.PostActData(ac)

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_CheckAP_PlaceCard, ac)
	})

	ac.PostActStream(func() {
		callback()
		return
	})
}
