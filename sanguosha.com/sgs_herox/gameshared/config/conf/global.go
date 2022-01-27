package conf

import "sanguosha.com/sgs_herox/proto/gameconf"

type Global struct {
	*gameconf.GlobalconfDefine
}

func (p *Global) loadConf(baseCfg *gameconf.GameBaseConfig) error {
	p.GlobalconfDefine = baseCfg.GetGlobalconf()[0]
	return nil
}
