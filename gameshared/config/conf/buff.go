package conf

import "sanguosha.com/sgs_herox/proto/gameconf"

type Buff struct {
	buffs map[int32]*gameconf.BuffConfDefine
}

func (p *Buff) loadConf(baseCfg *gameconf.GameBaseConfig) error {
	p.buffs = make(map[int32]*gameconf.BuffConfDefine)
	for _, v := range baseCfg.GetBuffConf() {
		p.buffs[v.BuffID] = v
	}
	return nil
}

func (p *Buff) Get(skillId int32) (*gameconf.BuffConfDefine, bool) {
	v, ok := p.buffs[skillId]
	return v, ok
}

func (p *Buff) All() map[int32]*gameconf.BuffConfDefine {
	return p.buffs
}
