package gameutil

import (
	"sort"
	// "github.com/sirupsen/logrus"
)

// 实现权重抽取

type DrawItem interface {
	GetID() int32
	GetWeight() int32
}

type DrawItemImp struct {
	id     int32
	weight int32
}

func NewDrawItem(id int32, weight int32) *DrawItemImp {
	return &DrawItemImp{
		id:     id,
		weight: weight,
	}
}

func (d *DrawItemImp) GetID() int32 {
	return d.id
}

func (d *DrawItemImp) GetWeight() int32 {
	return d.weight
}

type WeightDrawer struct {
	totalWeight  int32
	items        []DrawItem
	weightScales []int32
}

func NewWeightDrawer() *WeightDrawer {
	return &WeightDrawer{}
}

func (w WeightDrawer) Len() int {
	return len(w.items)
}

func (w *WeightDrawer) Swap(i, j int) {
	w.items[i], w.items[j] = w.items[j], w.items[i]
}

func (w *WeightDrawer) Less(i, j int) bool {
	if w.items[i].GetWeight() < w.items[j].GetWeight() {
		return true
	}
	return false
}

func (w *WeightDrawer) AddItem(items ...DrawItem) {
	for _, v := range items {
		if !w.HasItem(v.GetID()) {
			w.items = append(w.items, v)
		}
	}
	sort.Sort(w)
	w.weightScales = []int32{}
	w.totalWeight = 0
	for _, v := range w.items {
		if v.GetWeight() == 0 {
			continue
		}
		w.totalWeight += v.GetWeight()
		w.weightScales = append(w.weightScales, w.totalWeight)
	}
}

func (w *WeightDrawer) DrawOneItem() DrawItem {
	randNum := RandNum(w.totalWeight)
	for i, v := range w.weightScales {
		if randNum < v {
			return w.items[i]
		}
	}

	return nil
}

func (w *WeightDrawer) RemoveItem(id int32) {
	items := w.items
	w.items = []DrawItem{}
	w.totalWeight = 0
	w.weightScales = []int32{}
	for _, v := range items {
		if v.GetID() != id {
			w.AddItem(v)
		}
	}
}

func (w *WeightDrawer) IsEmpty() bool {
	return w.totalWeight == 0
}

func (w *WeightDrawer) HasItem(id int32) bool {
	for _, v := range w.items {
		if v.GetID() == id {
			return true
		}
	}
	return false
}

func (w *WeightDrawer) GetAll() []DrawItem {
	return w.items
}
