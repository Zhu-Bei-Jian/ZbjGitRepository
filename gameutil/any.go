package gameutil

import (
	"bytes"
	"encoding/json"
	"reflect"
	"runtime"
	"strconv"
)

type PackageInfo struct {
	Title string
	Info  []interface{}
}

func (a *PackageInfo) Append(info ...interface{}) *PackageInfo {
	a.Info = append(a.Info, info...)
	return a
}

func MakePackageInfo(title string, info ...interface{}) *PackageInfo {
	a := &PackageInfo{Title: title}
	return a.Append(info...)
}

func AnyIntValue(arg interface{}) int64 {
	switch arg.(type) {
	case int8:
		return int64(arg.(int8))
	case uint8:
		return int64(arg.(uint8))
	case int16:
		return int64(arg.(int16))
	case uint16:
		return int64(arg.(uint16))
	case int32:
		return int64(arg.(int32))
	case uint32:
		return int64(arg.(uint32))
	case int64:
		return int64(arg.(int64))
	case uint64:
		return int64(arg.(uint64))
	case int:
		return int64(arg.(int))
	case uint:
		return int64(arg.(uint))
	case string:
		v, _ := strconv.ParseInt(arg.(string), 10, 64)
		return v
	}
	return 0
}

func AnyStringValue(arg interface{}) string {
	switch arg.(type) {
	case int8:
		return strconv.Itoa(int(arg.(int8)))
	case uint8:
		return strconv.Itoa(int(arg.(uint8)))
	case int16:
		return strconv.Itoa(int(arg.(int16)))
	case uint16:
		return strconv.Itoa(int(arg.(uint16)))
	case int32:
		return strconv.Itoa(int(arg.(int32)))
	case uint32:
		return strconv.FormatInt(int64(arg.(uint32)), 10)
	case int64:
		return strconv.FormatInt(arg.(int64), 10)
	case uint64:
		return strconv.FormatUint(arg.(uint64), 10)
	case int:
		return strconv.Itoa(int(arg.(int)))
	case uint:
		return strconv.FormatInt(int64(arg.(uint)), 10)
	case string:
		return arg.(string)
	case []byte:
		return string(arg.([]byte))
	default:
		if arg != nil {
			bytes, e := json.Marshal(arg)
			if e == nil {
				return string(bytes)
			}
		}
	}
	return ""
}

func GetGoutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func NewStructMember(ptr interface{}, member string) interface{} {
	// 获取入参的类型
	t := reflect.TypeOf(ptr)
	// 入参类型校验
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil
	}
	// 取指针指向的结构体变量
	v := reflect.ValueOf(ptr).Elem()
	// 解析字段
	for i := 0; i < v.NumField(); i++ {
		// 取tag
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		// 解析tag
		memberName := tag.Get("memberName")
		if memberName != member {
			continue
		}
		if fieldInfo.Type.Kind() == reflect.Ptr {
			newValue := reflect.New(fieldInfo.Type.Elem())
			return newValue.Interface()
		}
		return reflect.New(fieldInfo.Type).Interface()
	}
	return nil
}
