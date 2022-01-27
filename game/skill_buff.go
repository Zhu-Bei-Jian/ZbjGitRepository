package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type ValueType int32

const (
	ValueType_HP ValueType = iota
	ValueType_Attack
	ValueType_Damage
	ValueType_ExpireV
)

type Buff interface {
	TriggerHandler() []TriggerHandler
	SetCfg(cfg *gameconf.BuffConfDefine)
	GetBuffId() int32
	OnEnable(card *Card)
	OnDisable(card *Card)
	GetValue(valueType ValueType) int32
	GetCfg() *gameconf.BuffConfDefine
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

func (p *HeroBuff) GetCfg() *gameconf.BuffConfDefine {
	return &p.cfg
}

func (p *HeroBuff) GetValue(valueType ValueType) int32 {
	switch valueType {
	case ValueType_HP:
		return p.cfg.BuffHP
	case ValueType_Attack:
		return p.cfg.BuffAttack
	case ValueType_Damage:
		return p.cfg.BuffDamage
	case ValueType_ExpireV:
		return p.cfg.ExpireV
	default:
		return -1
	}
}

func (p *HeroBuff) GetBuffId() int32 {
	return p.cfg.BuffID
}

func (p *HeroBuff) OnEnable(card *Card) {

}

func (p *HeroBuff) OnDisable(card *Card) {

}

func buffFeat(buffId int32) (Buff, bool) {
	switch buffId {
	case 0: //通用buff，用于清除其他buff
		return &BuffCommon{}, true
	case 1: //蔡文姬 断肠
		return &HeroBuff{}, true
	case 2: //曹操 奸雄
		return &BuffAtkAndHp{}, true
	case 3: //陈宫 智迟
		return &BuffZhiChi{}, true
	case 4: //大乔 国色
		return &BuffCommonMaxAttackCount{maxCount: INF}, true
	case 5: //法正
		return &HeroBuff{}, true
	case 6: //关羽
		return &HeroBuff{}, true
	case 7: //黄月英
		return &BuffAtkAndHp{}, true
	case 8: // 姜维
		return &HeroBuff{}, true
	case 9: // 刘备 激将
		return &BuffAtkAndHp{}, true
	case 10: //鲁肃 缔盟
		return &BuffDiMeng{}, true
	case 11: //马超
		return &BuffAtkAndHp{}, true
	case 12: //马谡
		return &BuffShiShou{}, true
	case 13: //孟获
		return &BuffAtkAndHp{}, true
	case 14: //糜竺
		return &BuffAtkAndHp{}, true
	case 15: //庞统
		return &BuffFengChu{}, true
	case 16: //司马徽
		return &HeroBuff{}, true
	case 17: //司马懿
		return &BuffAtkAndHp{}, true
	case 18: //孙坚
		return &BuffAtkAndHp{}, true
	case 19: //孙尚香
		return &BuffAtkAndHp{}, true
	case 20: //许褚 裸衣 +2
		return &BuffAtkAndHp{}, true
	case 21: //许褚 裸衣 +4
		return &BuffAtkAndHp{}, true
	case 22: //荀彧
		return &BuffJieMing{}, true
	case 23: //于禁
		return &HeroBuff{}, true
	case 24: //袁术
		return &BuffAtkAndHp{}, true
	case 25: //张郃
		return &BuffCommonAttackDistance{distance: INF}, true
	case 26: //张辽
		return &BuffAtkAndHp{}, true
	case 27: //赵云 龙胆
		return &HeroBuff{}, true
	case 28: //甄姬
		return &BuffAtkAndHp{}, true
	case 29: //孔融 礼让
		return &BuffAtkAndHp{}, true
	case 30: //曹仁
		return &BuffJuShou{}, true
	case 31: //甘夫人
		return &BuffAtkAndHp{}, true
	case 32: // 刘禅 放权
		return &BuffAtkAndHp{}, true
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

func newBuff(buffCfg *gameconf.BuffConfDefine) (Buff, bool) {
	buff1, ok := buffFeat(buffCfg.BuffID)
	if !ok {
		return nil, false
	}
	buff1.SetCfg(buffCfg)
	switch buff1.(type) {
	case *BuffAtkAndHp:
		bf := buff1.(*BuffAtkAndHp)
		bf.SetHP(buff1.GetCfg().GetBuffHP())
		bf.SetAttack(buff1.GetCfg().GetBuffAttack())
	default:
		//nothing
	}
	return buff1, ok
}
