package random

import (
	"errors"
	"sanguosha.com/sgs_herox/gameutil"
)

type group struct {
	totalRate int32
	items     []*groupItem

	readyPick bool
}

type groupItem struct {
	totalRate int32
	rate      int32
	item      interface{}
}

func NewGroup() *group {
	return &group{}
}

func (p *group) Add(rate int32, item interface{}) {
	if rate == 0 {
		return
	}
	p.items = append(p.items, &groupItem{
		totalRate: 0,
		rate:      rate,
		item:      item,
	})
}

func (p *group) checkPreparePick() {
	if p.readyPick {
		return
	}
	p.calculateRate()
	p.readyPick = true
}

func (p *group) calculateRate() {
	var totalRate int32
	for _, v := range p.items {
		totalRate += v.rate
		v.totalRate = totalRate
	}
	p.totalRate = totalRate
}

func (p *group) Pick() (interface{}, error) {
	p.checkPreparePick()

	if p.totalRate == 0 {
		return nil, errors.New("totalRate == 0")
	}
	randNum := gameutil.RandNum(p.totalRate)
	for _, v := range p.items {
		if v.totalRate > randNum {
			return v.item, nil
		}
	}
	return nil, errors.New("pick null")
}
