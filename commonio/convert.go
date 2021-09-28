package commonio

import (
	"reflect"
)

func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

func ConvertInterface2Int64(data interface{}) (int64, bool) {
	switch data.(type) {
	case float64:
		return int64(data.(float64)), true
	default:
		return 0, false
	}
}

func ConvertInterface2Int(data interface{}) (int, bool) {
	switch data.(type) {
	case float64:
		return int(data.(float64)), true
	default:
		return 0, false
	}
}

func ConvertInterface2String(data interface{}) (string, bool) {
	switch data.(type) {
	case string:
		return data.(string), true
	default:
		return "", false
	}
}
