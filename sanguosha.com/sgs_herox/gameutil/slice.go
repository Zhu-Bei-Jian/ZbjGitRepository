package gameutil

import (
	"sanguosha.com/sgs_herox/proto/db"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

func Shuffle(slice []interface{}) {
	for len(slice) > 0 {
		n := len(slice)
		randIndex := int(RandNum(int32(n)))
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func ShufflePropPack(slice []*gameconf.PropPack) {
	for len(slice) > 0 {
		n := len(slice)
		randIndex := int(RandNum(int32(n)))
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func ShuffleInt32(slice []int32) {
	for len(slice) > 0 {
		n := len(slice)
		randIndex := int(RandNum(int32(n)))
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func ShuffleInt(slice []int) {
	for len(slice) > 0 {
		n := len(slice)
		randIndex := int(RandNum(int32(n)))
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func IsInSlice(id int32, ids []int32) bool {
	for _, v := range ids {
		if id == v {
			return true
		}
	}
	return false
}

func IsInKVSliceK(key int32, kvs []*db.Int32KV) bool {
	for _, v := range kvs {
		if key == int32(v.Key) {
			return true
		}
	}
	return false
}

func IsInKVSliceV(value int32, kvs []*db.Int32KV) bool {
	for _, v := range kvs {
		if value == int32(v.V) {
			return true
		}
	}
	return false
}

func FindInSlice(id int32, ids []int32) int {
	for i, v := range ids {
		if id == v {
			return i
		}
	}
	return -1
}

func IsInStringSlice(id string, ids []string) bool {
	for _, v := range ids {
		if id == v {
			return true
		}
	}
	return false
}
