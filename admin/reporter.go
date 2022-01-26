package admin

import "sanguosha.com/sgs_herox/gameshared/reporter"

type ReporterManager struct {
	reporters []reporter.Reporter
}

func newReporterManager(tokens []string) (*ReporterManager, error) {
	var reporters []reporter.Reporter
	for _, v := range tokens {
		r, err := reporter.New(reporter.DingDing, v)
		if err != nil {
			return nil, err
		}
		reporters = append(reporters, r)
	}

	return &ReporterManager{
		reporters: reporters,
	}, nil
}

func (p *ReporterManager) Send(msg string) {
	for _, v := range p.reporters {
		v.Send(msg)
	}
}
