package conf

import "sanguosha.com/sgs_herox/proto/gameconf"

type Hero struct {
	heroes map[int32]*gameconf.HeroDefine
}

func (p *Hero) loadConf(baseCfg *gameconf.GameBaseConfig) error {
	p.heroes = make(map[int32]*gameconf.HeroDefine)
	for _, v := range baseCfg.GetHero() {
		p.heroes[v.HeroID] = v
	}
	return nil
}

func (p *Hero) GetHero(heroId int32) (*gameconf.HeroDefine, bool) {
	v, ok := p.heroes[heroId]
	return v, ok
}

func (p *Hero) HeroIds() (ret []int32) {
	for _, v := range p.heroes {
		ret = append(ret, v.HeroID)
	}
	return
}
func (p *Hero) GetHeroes() *map[int32]*gameconf.HeroDefine {
	return &p.heroes
}
