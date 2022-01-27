package gameutil

import (
	"math"
	"strconv"
	"strings"
)

func SafeAdd(p interface{}, v int64) {
	if v <= 0 {
		return
	}
	switch p.(type) {
	case *int64:
		add := v
		pValue := p.(*int64)
		if *pValue < math.MaxInt64 {
			less := math.MaxInt64 - *pValue
			if add > less {
				add = less
			}
			*pValue += add
		}
	case *int32:
		if v > int64(math.MaxInt32) {
			return
		}
		add := int32(v)
		pValue := p.(*int32)
		if *pValue < math.MaxInt32 {
			less := math.MaxInt32 - *pValue
			if add > less {
				add = less
			}
			*pValue += add
		}
	}
}

func MinAndDiff(x int32, y int32) (int32, int32) {
	if x < y {
		return x, y - x
	}
	return y, x - y
}

func Min(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func MinInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x int32, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func Abs(x int32) int32 {
	if x > 0 {
		return x
	} else {
		return -x
	}
}

func Int64Min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Int64Max(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Int32Min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func Int32Max(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

//生成不重复的随机数，随机数范围0~max,数量num
func GenDiffRandomNum(num, max int) (outIntArr []int) {
	if num >= max {
		for i := 0; i < int(max); i++ {
			outIntArr = append(outIntArr, i)
		}
		return
	}
	var tmpArr []int
	for i := 0; i < int(max); i++ {
		tmpArr = append(tmpArr, 1)
	}

	var rnd int

	for i := 0; i < int(num); i++ {
		for {
			rnd = int(RandNum(int32(max)))
			if tmpArr[rnd] != -1 {
				break
			}
		}
		outIntArr = append(outIntArr, rnd)
		tmpArr[rnd] = -1
	}
	return
}

func ToInt32(s string) (int32, error) {
	i, err := strconv.Atoi(strings.TrimSpace(s))
	return int32(i), err
}

func ToInt(s string) (int, error) {
	i, err := strconv.Atoi(strings.TrimSpace(s))
	return i, err
}

func ToBool(s string) (bool, error) {
	v, err := strconv.ParseBool(s)
	return v, err
}

func ToUInt64(s string) (uint64, error) {
	i, err := strconv.ParseUint(strings.TrimSpace(s), 0, 64)
	return uint64(i), err
}

func ToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(strings.TrimSpace(s), 0, 64)
	return i, err
}

func ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func ToFloat64(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	return v, err
}

//取最后两位数
func Last2Num(num int32) int32 {
	return num - 100*(num/100)
}

func LastNum(num int32) int32 {
	return num - 10*(num/10)
}

//int32是否溢出
func IsOverflowInt32(v int64) bool {
	return v > int64(math.MaxInt32)
}

//安全相乘
func SafeMultiInt32(v1, v2 int32) int64 {
	return int64(v1) * int64(v2)
}

func IsNowValid(startTime, endTime int64) bool {
	now := GetCurrentTimestamp()
	if startTime != 0 && now < startTime {
		return false
	}
	if endTime != 0 && now > endTime {
		return false
	}
	return true
}
