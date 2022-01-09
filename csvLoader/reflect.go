package csv

import (
	"reflect"
	"strings"
)

// FieldKind 数据类型。 注:为什么不直接用 reflect.Kind, 因为不是所有类型都支持转化，不支持的自动转变为 Invalid
type FieldKind uint8
const (
	Invalid FieldKind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	Map
	Slice
	String
)

type fieldInfo struct {
	name    string
	kind    FieldKind
	index   int
	context reflect.Type
}

func getFieldKind(t reflect.Kind) FieldKind {
	switch t {
	case reflect.String:
		return String
	case reflect.Bool:
		return Bool
	case reflect.Int:
		return Int
	case reflect.Int8:
		return Int8
	case reflect.Int16:
		return Int16
	case reflect.Int32:
		return Int32
	case reflect.Int64:
		return Int64
	case reflect.Uint:
		return Uint
	case reflect.Uint8:
		return Uint8
	case reflect.Uint16:
		return Uint16
	case reflect.Uint32:
		return Uint32
	case reflect.Uint64:
		return Uint64
	case reflect.Float32:
		return Float32
	case reflect.Float64:
		return Float64
	case reflect.Map:
		return Map
	case reflect.Slice:
		return Slice
	default:
		return Invalid
	}
}

func getFieldIndex(name string, data *[]string) int {
	realName := strings.ToLower(name)
	for i, v := range *data {
		if v == realName {
			return i
		}
	}

	return -1
}

func setFieldValue(v interface{}, infoList *[]*fieldInfo, data *[]string) {
	vValue := reflect.ValueOf(v).Elem()
	for i, info := range *infoList {
		if info.index > -1 {
			value, err := fieldParser[info.kind]((*data)[info.index], info.context)
			if err == nil {
				vValue.Field(i).Set(value)
			}
		}
	}
}
