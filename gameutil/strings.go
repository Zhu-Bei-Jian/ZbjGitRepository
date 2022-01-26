package gameutil

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ToLowerTrim(s string) string {
	// 过滤前后空格并转为小写
	return strings.ToLower(strings.TrimSpace(s))
}

func RandString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	tmp := base64.URLEncoding.EncodeToString(b)
	return string([]byte(tmp)[:n])
}

func GetRandomString(length int) string {
	//str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomStringAll(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func MakeTokenV1(id uint64, t uint64) string {
	randKey := uint64(RandNum(8000)) + 1
	token := fmt.Sprintf("%d-%d-%d", 1357+randKey, (id/13*randKey)%10000, SafeDivUint64(t, (57*randKey))%10000)
	return token
}
func CheckTokenV1(id uint64, token string) bool {
	list, _ := ParseUint64s(token, "-")
	if len(list) != 3 {
		return false
	}
	randKey := list[0] - 1357
	if list[1] != (id/13*uint64(randKey))%10000 {
		return false
	}
	return true
}

func ParseUInt32s(s string, sep string) ([]uint32, error) {
	ids := strings.Split(s, sep)
	res := make([]uint32, 0)
	for _, v := range ids {
		id, err := ToInt32(v)
		if err != nil {
			return nil, err
		}
		res = append(res, uint32(id))
	}
	return res, nil
}

func ParseInt32s(s string, sep string) ([]int32, error) {
	ids := strings.Split(s, sep)
	res := make([]int32, 0)
	for _, v := range ids {
		id, err := ToInt32(v)
		if err != nil {
			return nil, err
		}
		res = append(res, int32(id))
	}
	return res, nil
}

func ParseUint64s(s string, sep string) ([]uint64, error) {
	ids := strings.Split(s, sep)
	res := make([]uint64, 0)
	for _, v := range ids {
		id, err := ToUInt64(v)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func StringToInt32List(s string, sep string) []int32 {
	ids := strings.Split(s, sep)
	res := make([]int32, 0)
	for _, v := range ids {
		id, err := ToInt32(v)
		if err != nil {
			return res
		}
		res = append(res, int32(id))
	}
	return res
}

func ParseDateTimes(s string) ([]int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	times := strings.Split(s, ";")
	res := make([]int64, 0)
	for _, v := range times {
		id, err := ParseDatetime2Timestamp(v)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func GetPayHash(params map[string]interface{}, secretKey string) string {
	// 排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 生成sign
	signData := ""
	for _, v := range keys {
		signData += fmt.Sprintf("%s=%v", v, params[v])
	}
	// 生成签名
	return SHA512(signData + secretKey)
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------
// 名字中的表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}
