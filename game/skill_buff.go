package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type Buff interface {
	TriggerHandler() []TriggerHandler
	SetCfg(cfg *gameconf.BuffConfDefine)
	GetBuffId() int32
	OnEnable(card *Card)
	OnDisable(card *Card)
}

type HeroBuff struct {
	Buff
	core.ActionDataCore
	cfg gameconf.BuffConfDefine
}

func (p *HeroBuff) TriggerHandler() []TriggerHandler {
	return []TriggerHandler{}
}

func (p *HeroBuff) SetCfg(cfg *gameconf.BuffConfDefine) {
	p.cfg = *cfg
}

func (p *HeroBuff) GetBuffId() int32 {
	return p.cfg.BuffID
}

func (p *HeroBuff) OnEnable(card *Card) {

}

func (p *HeroBuff) OnDisable(card *Card) {

}

func NewBuff(buffCfg *gameconf.BuffConfDefine) (Buff, bool) {
	buff1, ok := newBuff(buffCfg.BuffID)
	if !ok {
		return nil, false
	}
	buff1.SetCfg(buffCfg)
	return buff1, true
}

func newBuff(buffId int32) (Buff, bool) {
	switch buffId {
	case 0: //通用buff，用于清除其他buff
		return &BuffCommon{}, true
	case 1: //蔡文姬 断肠
		return &HeroBuff{}, true
	case 2: //曹操 奸雄
		return &HeroBuff{}, true
	case 3: //陈宫 智迟
		return &BuffZhiChi{}, true
	case 4: //大乔 国色
		return &BuffCommonMaxAttackCount{maxCount: INF}, true
	case 5: //法正
		return &HeroBuff{}, true
	case 6: //关羽
		return &HeroBuff{}, true
	case 7: //黄月英
		return &HeroBuff{}, true
	case 8: // 姜维
		return &HeroBuff{}, true
	case 9: // 刘备 激将
		return &BuffAtkAndHp{attack: 2}, true
	case 10: //鲁肃 缔盟
		return &BuffDiMeng{}, true
	case 11: //马超
		return &BuffAtkAndHp{attack: 10}, true
	case 12: //马谡
		return &BuffShiShou{}, true
	case 13: //孟获
		return &HeroBuff{}, true
	case 14: //糜竺
		return &HeroBuff{}, true
	case 15: //庞统
		return &BuffFengChu{}, true
	case 16: //司马徽
		return &HeroBuff{}, true
	case 17: //司马懿
		return &HeroBuff{}, true
	case 18: //孙坚
		return &HeroBuff{}, true
	case 19: //孙尚香
		return &BuffAtkAndHp{attack: 1}, true
	case 20: //许褚 裸衣 +2
		return &BuffAtkAndHp{attack: 2}, true
	case 21: //许褚 裸衣 +4
		return &BuffAtkAndHp{attack: 4}, true
	case 22: //荀彧
		return &BuffJieMing{}, true
	case 23: //于禁
		return &HeroBuff{}, true
	case 24: //袁术
		return &HeroBuff{}, true
	case 25: //张郃
		return &BuffCommonAttackDistance{distance: INF}, true
	case 26: //张辽
		return &BuffAtkAndHp{attack: 5}, true
	case 27: //赵云 龙胆
		return &HeroBuff{}, true
	case 28: //甄姬
		return &BuffAtkAndHp{attack: 3}, true
	case 29: //孔融 礼让
		return &BuffAtkAndHp{hp: 3}, true
	case 30: //曹仁
		return &BuffJuShou{}, true
	case 31: //甘夫人
		return &BuffAtkAndHp{hp: 2}, true
	case 32: // 刘禅 放权
		return &BuffAtkAndHp{attack: 3, hp: 3}, true
	case 33: // 鲍三娘 许身
		return &BuffXuShen{}, true
	case 34: // 陈琳 颂词
		return &BuffSongCi{}, true
	case 35: // 曹丕 行殇
		return &BuffXingShang{}, true
	case 36:
		return &BuffZiShou{}, true
	//case
	default:
		return nil, false
	}
}
