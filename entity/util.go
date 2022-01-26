package entity

import "sort"

func IsUniqueCardGroup(cardGroup []int32) bool {
	sort.Slice(cardGroup, func(i, j int) bool {
		return cardGroup[i] < cardGroup[j]
	})
	for i := 1; i < len(cardGroup); i++ {
		if cardGroup[i] == cardGroup[i-1] {
			return false
		}
	}
	return true
}
