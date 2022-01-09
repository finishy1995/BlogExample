package csv

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type testData struct {
	A string
	B int32
	D *bool   // 错误的数据类型
	E string // 不存在的
}

func TestLoadDataFromCsvBytes(t *testing.T) {
	r := require.New(t)
	in := `---,A,C,B,D
8,"aa",11,12,true
10,"bb",21,22,true
1000,"cc",21,22,false
`
	nullIn := ""
	wrongIn := `---,A,C,B
8,"aa",11,12
10,"bb",21,22
1000,"cc",21
`
	io := strings.NewReader(in)
	nullIO := strings.NewReader(nullIn)
	wrongIO := strings.NewReader(wrongIn)

	// 测试错误的输入类型
	testNumber := 1
	result, err := loadDataFromCsvBytes(io, testNumber)
	r.Equal(ErrInvalidInputStruct, err)
	result, err = loadDataFromCsvBytes(io, &testNumber)
	r.Equal(ErrInvalidInputStruct, err)

	// 测试错误的csv类型
	var testStruct testData
	result, err = loadDataFromCsvBytes(nullIO, &testStruct)
	r.Nil(err)
	r.Nil(result)
	result, err = loadDataFromCsvBytes(wrongIO, &testStruct)
	r.NotNil(err)

	// 测试正确情况
	result, err = loadDataFromCsvBytes(io, &testStruct)
	r.Nil(err)
	r.Nil(result["1"])
	r.NotNil(result["10"])
	s := result["10"].(*testData)
	r.Equal("bb", s.A)
	r.Equal(int32(22), s.B)
	var ptr *bool
	r.Equal(ptr, s.D)
	r.Equal("", s.E)
}

type testAllType struct {
	A string
	B bool
	C int
	D int8
	E int16
	F int32
	G int64
	H uint
	I uint8
	J uint16
	K uint32
	L uint64
	M float32
	N float64
	O map[string]bool
	P []uint32
}

func TestAllDataType(t *testing.T) {
	r := require.New(t)
	in := `---,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P
1,aa,true,"1",-2,-3,4,-5,6,7,8,9,10,11,12.9,"(a:false,""b"":true,""c"":aaa)","(111,-222,11)"
`
	io := strings.NewReader(in)
	var testStruct testAllType
	result, err := loadDataFromCsvBytes(io, &testStruct)
	r.Nil(err)
	r.Equal(1, len(result))
	item := result["1"].(*testAllType)
	r.NotNil(item)
	r.Equal("aa", item.A)
	r.Equal(true, item.B)
	r.Equal(1, item.C)
	r.Equal(int8(-2), item.D)
	r.Equal(int16(-3), item.E)
	r.Equal(int32(4), item.F)
	r.Equal(int64(-5), item.G)
	r.Equal(uint(6), item.H)
	r.Equal(uint8(7), item.I)
	r.Equal(uint16(8), item.J)
	r.Equal(uint32(9), item.K)
	r.Equal(uint64(10), item.L)
	r.Equal(float32(11), item.M)
	r.Equal(12.9, item.N)
	r.Equal(2,len(item.O))
	r.Equal(false, item.O["a"])
	r.Equal(true, item.O["b"])
	r.Equal(2,len(item.P))
	r.Equal(uint32(111), item.P[0])
	r.Equal(uint32(11), item.P[1])
}

func TestLoadData(t *testing.T) {
	r := require.New(t)
	var testStruct testData
	_, err := LoadData("test/data.csv",&testStruct)
	r.Nil(err)
	_, err = LoadData("test/data.t",&testStruct)
	r.NotNil(err)
	_, err = LoadData("test/data1.csv",&testStruct)
	r.NotNil(err)
}

func TestMustLoad(t *testing.T) {
	r := require.New(t)
	var testStruct testData
	r.NotPanics(func() {
		MustLoad("test/data.csv",&testStruct)
	})
	r.Panics(func() {
		MustLoad("test/data1.csv",&testStruct)
	})
}
