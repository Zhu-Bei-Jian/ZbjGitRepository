package gameutil

import (
	"math/rand"
	"sync"
	"time"
)

//var myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
//保证多线程访问安全
var myRand = rand.New(&LockdSource{src: rand.NewSource(time.Now().UnixNano()).(rand.Source64)})

type LockdSource struct {
	lk  sync.Mutex
	src rand.Source64
}

func (r *LockdSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *LockdSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

// RandNum ...
func RandNum(num int32) int32 {
	return myRand.Int31n(num)
}

func Rand() int {
	return myRand.Int()
}

func Rand63n(v int64) int64 {
	return myRand.Int63n(v)
}

func Perm(l int) []int {
	return myRand.Perm(l)
}
func Rand31() int32 {
	return myRand.Int31()
}

func Rand63() int64 {
	return myRand.Int63()
}

func RandPickListInt32(data []int32, num int) (r []int32) { //ReservoirSampling
	total := len(data)
	for i := 0; i < total && i < num; i++ {
		r = append(r, data[i])
	}
	for i := num; i < total; i++ {
		tmp := myRand.Intn(i + 1)
		if tmp < num {
			r[tmp] = data[i]
		}
	}
	return
}
func RandPickListInt32V2(data []int32, num int) (r []int32) {
	total := len(data)
	for i := 0; i < total && num > 0; i++ {
		tmp := myRand.Intn(total - i)
		if tmp < num {
			r = append(r, data[i])
			num--
		}
	}
	return
}

func RandPickIndex(total int, num int) (r []int) { //ReservoirSampling
	for i := 0; i < total && i < num; i++ {
		r = append(r, i)
	}
	for i := num; i < total; i++ {
		tmp := myRand.Intn(i + 1)
		if tmp < num {
			r[tmp] = i
		}
	}
	return
}

func RandListInt32(data []int32) {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}
func RandListUInt64(data []uint64) {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

// RandomBetween ...
func RandomBetween(x, y int32) int32 {
	d := x - y
	if d == 0 {
		return x
	} else if d < 0 {
		return x + RandNum(1-d)
	} else {
		return y + RandNum(d+1)
	}
}

// RandomBetweenZero ...
func RandomBetweenZero(min, max int) int {
	if min >= max {
		return max
	}
	return myRand.Intn(max-min+1) + min
}

// RandomBetween31n ...
func RandomBetween31n(min, max int32) int32 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return myRand.Int31n(max-min+1) + min
}

// IsProbIntn ...
func IsProbIntn(percent int) bool {
	return myRand.Intn(100) < percent
}

// IsProbInt31n ...
func IsProbInt31n(percent int32) bool {
	return myRand.Int31n(100) < percent
}

// GetRandomN ...
func GetRandomN(total, num int) []int {
	if total < num {
		num = total
	}
	idx := myRand.Perm(total)
	res := make([]int, 0, num)
	for i := 0; i < num; i++ {
		res = append(res, idx[i])
	}
	return res
}

func GetRandomListNoRepeatOnCheck(pool []int32, n int32, check func(int32) bool, idMap map[int32]int32, limit int32) (r []int32) {
	tmp := append([]int32{}, pool...)
	left := int32(len(tmp))
	for n > 0 && left > 0 {
		i := RandNum(left)
		v := tmp[i]
		if idMap != nil && limit != 0 {
			if n, ok := idMap[v]; ok {
				if n >= limit {
					tmp = append(tmp[0:i], tmp[i+1:]...)
					left--
					//logrus.Debug("GetRandomNWithoutRepeatOfLimit skip:", v)
					continue
				} else {
					idMap[v] = n + 1
				}
			} else {
				idMap[v] = 1
			}
		}
		if check == nil || check(v) {
			r = append(r, v)
			n--
		}
		tmp = append(tmp[0:i], tmp[i+1:]...)
		left--
	}
	return
}

func GetRandomNoRepeat(pool []int32, n int32) (r []int32) {
	tmp := append([]int32{}, pool...)
	left := int32(len(tmp))
	for n > 0 && left > 0 {
		i := RandNum(left)
		r = append(r, tmp[i])
		n--
		tmp = append(tmp[0:i], tmp[i+1:]...)
		left--
	}
	return
}

func GetRandomNWithoutRepeatOfLimit(pool []int32, n int32, idMap map[int32]int32, limit int32) (r []int32) {
	tmp := append([]int32{}, pool...)
	left := int32(len(tmp))
	for n > 0 && left > 0 {
		i := RandNum(left)
		v := tmp[i]
		if idMap != nil && limit != 0 {
			if n, ok := idMap[v]; ok {
				if n >= limit {
					tmp = append(tmp[0:i], tmp[i+1:]...)
					left--
					//logrus.Debug("GetRandomNWithoutRepeatOfLimit skip:", v)
					continue
				} else {
					idMap[v] = n + 1
				}
			} else {
				idMap[v] = 1
			}
		}
		r = append(r, tmp[i])
		n--
		tmp = append(tmp[0:i], tmp[i+1:]...)
		left--
	}
	//logrus.Debug("GetRandomNWithoutRepeatOfLimit", pool, r)
	return
}
