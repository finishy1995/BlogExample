package csv

import (
	"reflect"
	"strconv"
	"strings"
)

func allToLower(content *[]string) {
	for i, v := range *content {
		(*content)[i] = strings.ToLower(v)
	}
}

type fieldHandler func(str string, context reflect.Type) (reflect.Value, error)
var (
	fieldParser = map[FieldKind]fieldHandler {
		Invalid: stringToInvalid,
		String: stringToString,
		Bool: stringToBool,
		Int: stringToInt,
		Int8: stringToInt8,
		Int16: stringToInt16,
		Int32: stringToInt32,
		Int64: stringToInt64,
		Uint: stringToUint,
		Uint8: stringToUint8,
		Uint16: stringToUint16,
		Uint32: stringToUint32,
		Uint64: stringToUint64,
		Float32: stringToFloat32,
		Float64: stringToFloat64,
		Map: stringToMap,
		Slice: stringToSlice,
	}
	fieldParserHelper map[FieldKind]fieldHandler
)

func init() {
	fieldParserHelper = fieldParser
}

func stringToInvalid(_ string, _ reflect.Type) (reflect.Value, error) {
	return reflect.Value{}, ErrUnsupportedDataType
}

func stringToString(str string, _ reflect.Type) (reflect.Value, error) {
	str = strings.Trim(str, "\"")
	return reflect.ValueOf(str), nil
}

func stringToBool(str string, _ reflect.Type) (reflect.Value, error) {
	realStr := strings.ToLower(str)
	if realStr == "true" || realStr == "y" || realStr == "yes" {
		return reflect.ValueOf(true), nil
	}
	if realStr == "false" || realStr == "n" || realStr == "no" {
		return reflect.ValueOf(false), nil
	}
	return reflect.Value{}, ErrInvalidDataSource
}

func stringToInt(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseInt(str, 10, 0)
	return reflect.ValueOf(int(result)), err
}

func stringToInt8(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseInt(str, 10, 8)
	return reflect.ValueOf(int8(result)), err
}

func stringToInt16(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseInt(str, 10, 16)
	return reflect.ValueOf(int16(result)), err
}

func stringToInt32(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseInt(str, 10, 32)
	return reflect.ValueOf(int32(result)), err
}

func stringToInt64(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseInt(str, 10, 64)
	return reflect.ValueOf(result), err
}

func stringToUint(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseUint(str, 10, 0)
	return reflect.ValueOf(uint(result)), err
}

func stringToUint8(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseUint(str, 10, 8)
	return reflect.ValueOf(uint8(result)), err
}

func stringToUint16(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseUint(str, 10, 16)
	return reflect.ValueOf(uint16(result)), err
}

func stringToUint32(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseUint(str, 10, 32)
	return reflect.ValueOf(uint32(result)), err
}

func stringToUint64(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseUint(str, 10, 64)
	return reflect.ValueOf(result), err
}

func stringToFloat32(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseFloat(str, 32)
	return reflect.ValueOf(float32(result)), err
}

func stringToFloat64(str string, _ reflect.Type) (reflect.Value, error) {
	result, err := strconv.ParseFloat(str, 64)
	return reflect.ValueOf(result), err
}

// stringToMap 格式 (A:"aaa",B:"ccc")，支持去除任何前后置空格
func stringToMap(str string, context reflect.Type) (reflect.Value, error) {
	if fieldParserHelper == nil {
		return reflect.Value{}, ErrInvalidDataSource
	}
	str = strings.Trim(str, " ")
	str = str[1:len(str)-1]
	strArr := strings.Split(str, ",")
	result := reflect.MakeMap(context)
	for _, s := range strArr {
		s = strings.Trim(s, " ")
		kv := strings.Split(s, ":")
		if len(kv) != 2 {
			continue
		}

		// 避免循环引用，顺便检查是否是支持的数据类型
		keyHandler, ok := fieldParserHelper[getFieldKind(context.Key().Kind())]
		if !ok {
			continue
		}
		valueKind, ok := fieldParserHelper[getFieldKind(context.Elem().Kind())]
		if !ok {
			continue
		}

		// 赋值
		key, err := keyHandler(kv[0], context.Key())
		if err != nil {
			continue
		}
		value, err := valueKind(kv[1], context.Elem())
		if err != nil {
			continue
		}
		result.SetMapIndex(key, value)
	}
	return result, nil
}

// stringToSlice 格式 ("aaa","ccc")，支持去除任何前后置空格
func stringToSlice(str string, context reflect.Type) (reflect.Value, error) {
	if fieldParserHelper == nil {
		return reflect.Value{}, ErrInvalidDataSource
	}
	str = strings.Trim(str, " ")
	str = str[1:len(str)-1]
	strArr := strings.Split(str, ",")
	result := reflect.MakeSlice(context, 0, len(strArr))
	for _, s := range strArr {
		s = strings.Trim(s, " ")
		if handler, ok := fieldParserHelper[getFieldKind(context.Elem().Kind())]; ok {
			value, err := handler(s, context.Elem())
			if err != nil {
				continue
			}
			result = reflect.Append(result, value)
		}
	}
	return result, nil
}
