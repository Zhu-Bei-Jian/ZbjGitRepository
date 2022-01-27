package gameutil

import "reflect"

func StructJsonToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Tag.Get("json")] = obj2.Field(i).Interface()
	}
	return data
}
