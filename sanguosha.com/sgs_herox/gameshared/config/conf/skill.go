package conf

import "sanguosha.com/sgs_herox/proto/gameconf"

type Skill struct {
	skills map[int32]*gameconf.SkillConfDefine
}

func (p *Skill) loadConf(baseCfg *gameconf.GameBaseConfig) error {
	p.skills = make(map[int32]*gameconf.SkillConfDefine)
	for _, v := range baseCfg.GetSkillConf() {
		p.skills[v.SkillID] = v
	}
	return nil
}

func (p *Skill) GetSkill(skillId int32) (*gameconf.SkillConfDefine, bool) {
	v, ok := p.skills[skillId]
	return v, ok
}

func (p *Skill) All(typ gameconf.SkillTyp) map[int32]*gameconf.SkillConfDefine {
	m := make(map[int32]*gameconf.SkillConfDefine)
	for _, v := range p.skills {
		if v.SkillType == typ {
			m[v.SkillID] = v
		}
	}
	return m
}
