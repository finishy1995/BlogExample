package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
)

var (
	loaders = map[string]func(io.Reader, interface{}) (map[string]interface{}, error){
		".csv": loadDataFromCsvBytes,
	}
	ErrInvalidInputStruct  = errors.New("invalid data input struct, make sure interface{} is *struct")
	ErrInvalidDataSource   = errors.New("invalid data source, please check csv data")
	ErrUnsupportedDataType = errors.New("unsupported data type")
)

// MustLoad 必须加载成功数据才能启动
func MustLoad(file string, v interface{}) map[string]interface{} {
	result, err := LoadData(file, v)
	if err != nil {
		panic(errors.New(fmt.Sprintf("error: data file %s, %s", file, err.Error())))
	}
	return result
}

// LoadData 加载数据，失败返回 error
func LoadData(file string, v interface{}) (map[string]interface{}, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	loader, ok := loaders[path.Ext(file)]
	if !ok {
		return nil, fmt.Errorf("unrecognized file type: %s", file)
	}

	return loader(f, v)
}

// loadDataFromCsvBytes loads data into v from content csv bytes.
func loadDataFromCsvBytes(content io.Reader, v interface{}) (map[string]interface{}, error) {
	// 检查元素类型
	if reflect.TypeOf(v).Kind() != reflect.Ptr || reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
		return nil, ErrInvalidInputStruct
	}

	// 读取数据
	reader := csv.NewReader(content)
	preData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// 检查处理表头
	length := len(preData) - 1
	if length < 0 {
		return nil, nil
	}
	allToLower(&preData[0])

	// 提取结构体参数名、类型、对应的位置
	vValue := reflect.ValueOf(v).Elem()
	fieldNum := vValue.NumField()
	vType := reflect.TypeOf(v).Elem()
	fieldInfoList := make([]*fieldInfo, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := vType.Field(i)
		fieldInfoList[i] = &fieldInfo{
			name:  field.Name,
			kind:  getFieldKind(field.Type.Kind()),
			index: getFieldIndex(field.Name, &preData[0]),
		}
		if fieldInfoList[i].kind == Map || fieldInfoList[i].kind == Slice {
			fieldInfoList[i].context = field.Type
		}
	}

	// 生成数据
	data := make(map[string]interface{}, length)
	for i := 0; i < length; i++ {
		data[preData[i+1][0]] = reflect.New(vValue.Type()).Interface()
		setFieldValue(data[preData[i+1][0]], &fieldInfoList, &preData[i+1])
	}

	return data, nil
}
