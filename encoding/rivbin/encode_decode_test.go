package rivbin

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"testing"
)

type ( // types to ensure our integral types can work with type aliases
	u   uint
	u8  uint8
	u16 uint16
	u32 uint32
	u64 uint64

	i   int
	i8  int8
	i16 int16
	i32 int32
	i64 int64
)

func TestMarshalUnmarshalDynamicIntegers(t *testing.T) {
	testCases := []struct {
		Value      interface{}
		ByteLength int
	}{
		// uint8
		{
			uint8(0),
			1,
		},
		{
			uint8(1),
			1,
		},
		{
			uint8(math.MaxUint8),
			1,
		},
		{
			u8(1),
			1,
		},
		// uint16
		{
			uint16(0),
			2,
		},
		{
			uint16(1),
			2,
		},
		{
			uint16(math.MaxUint16),
			2,
		},
		{
			u16(1),
			2,
		},
		// uint32
		{
			uint32(0),
			4,
		},
		{
			uint32(1),
			4,
		},
		{
			uint32(math.MaxUint32),
			4,
		},
		{
			u32(1),
			4,
		},
		// uint64
		{
			uint64(0),
			8,
		},
		{
			uint64(1),
			8,
		},
		{
			uint64(math.MaxUint64),
			8,
		},
		{
			u64(1),
			8,
		},
		// uint
		{
			uint(0),
			8,
		},
		{
			uint(1),
			8,
		},
		{
			uint(math.MaxUint32),
			8,
		},
		{
			uint(math.MaxUint64),
			8,
		},
		{
			u(1),
			8,
		},
		// int8
		{
			int8(0),
			1,
		},
		{
			int8(1),
			1,
		},
		{
			int8(math.MaxInt8),
			1,
		},
		{
			i8(1),
			1,
		},
		// int16
		{
			int16(0),
			2,
		},
		{
			int16(1),
			2,
		},
		{
			int16(math.MaxInt16),
			2,
		},
		{
			i16(1),
			2,
		},
		// int32
		{
			int32(0),
			4,
		},
		{
			int32(1),
			4,
		},
		{
			int32(math.MaxInt32),
			4,
		},
		{
			i32(1),
			4,
		},
		// int64
		{
			int64(0),
			8,
		},
		{
			int64(1),
			8,
		},
		{
			int64(math.MaxInt64),
			8,
		},
		{
			i64(1),
			8,
		},
		// int
		{
			int(0),
			8,
		},
		{
			int(1),
			8,
		},
		{
			int(math.MaxInt32),
			8,
		},
		{
			int(math.MaxInt64),
			8,
		},
		{
			i(1),
			8,
		},
	}
	for idx, testCase := range testCases {
		b := bytes.NewBuffer(nil)

		// marshal
		NewEncoder(b).Encode(testCase.Value)

		// test byte value
		bl := len(b.Bytes())
		if bl != testCase.ByteLength {
			t.Errorf("unexpected byte length of %d for testcase %d ('%v')", bl, idx, testCase.Value)
			// do continue, as to also test the other required properties,
			// an error here does not have to mean the other properties are broken as well
		}

		// write extra content, as to ensure nothing more is read, than should be read
		n, err := b.Write([]byte("test"))
		if err != nil {
			t.Fatal(idx, err)
		}
		if n != 4 {
			t.Fatal(idx, io.ErrShortWrite)
		}

		// allocate new value
		val := reflect.New(reflect.TypeOf(testCase.Value))

		// unmarshal (valid values only)
		err = NewDecoder(b).Decode(val.Interface())
		if err != nil {
			t.Errorf("unexpected Unmarshal error for testCase %d (iptrVal: %s/%s) ('%v'): %v",
				idx, val.Type().String(), val.Kind().String(), testCase.Value, err)
			continue
		}

		// compare unmarshalled value with original value
		if ival := val.Elem().Interface(); !reflect.DeepEqual(testCase.Value, ival) {
			t.Errorf("unexpected unmarshalled value: '%v' != '%v'", testCase.Value, ival)
			continue
		}

		// ensure that the unmarshal process has read the exact amount of bytes
		remaining := string(b.Bytes())
		if remaining != "test" {
			t.Error(idx, "read more or less than it should have, unexpected remaining:", remaining, "!= test")
		}
	}
}

func TestUnmarshalDynamicIntegersWithInvalidInput(t *testing.T) {
	testCases := []struct {
		Input       []byte
		OutputValue interface{}
	}{
		// uint8
		{
			nil,
			uint8(0),
		},
		{
			[]byte{},
			uint8(0),
		},
		// uint16
		{
			nil,
			uint16(0),
		},
		{
			[]byte{},
			uint16(0),
		},
		{
			[]byte{0},
			uint16(0),
		},
		{
			[]byte{1},
			uint16(0),
		},
		// uint32
		{
			nil,
			uint32(0),
		},
		{
			[]byte{},
			uint32(0),
		},
		{
			[]byte{0},
			uint32(0),
		},
		{
			[]byte{1},
			uint32(0),
		},
		{
			[]byte{1, 0},
			uint32(0),
		},
		{
			[]byte{1, 1},
			uint32(0),
		},
		{
			[]byte{1, 0, 1},
			uint32(0),
		},
		{
			[]byte{1, 1, 1},
			uint32(0),
		},
		// uint64
		{
			nil,
			uint64(0),
		},
		{
			[]byte{},
			uint64(0),
		},
		{
			[]byte{1},
			uint64(0),
		},
		{
			[]byte{1, 1},
			uint64(0),
		},
		{
			[]byte{1, 1, 1},
			uint64(0),
		},
		{
			[]byte{1, 1, 1, 1},
			uint64(0),
		},
		{
			[]byte{1, 1, 1, 1, 1},
			uint64(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1},
			uint64(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1, 1},
			uint64(0),
		},
		// uint
		{
			nil,
			uint(0),
		},
		{
			[]byte{},
			uint(0),
		},
		{
			[]byte{1},
			uint(0),
		},
		{
			[]byte{1, 1},
			uint(0),
		},
		{
			[]byte{1, 1, 1},
			uint(0),
		},
		{
			[]byte{1, 1, 1, 1},
			uint(0),
		},
		{
			[]byte{1, 1, 1, 1, 1},
			uint(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1},
			uint(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1, 1},
			uint(0),
		},
		// int8
		{
			nil,
			int8(0),
		},
		{
			[]byte{},
			int8(0),
		},
		// int16
		{
			nil,
			int16(0),
		},
		{
			[]byte{},
			int16(0),
		},
		{
			[]byte{0},
			int16(0),
		},
		{
			[]byte{1},
			int16(0),
		},
		// int32
		{
			nil,
			int32(0),
		},
		{
			[]byte{},
			int32(0),
		},
		{
			[]byte{0},
			int32(0),
		},
		{
			[]byte{1},
			int32(0),
		},
		{
			[]byte{1, 0},
			int32(0),
		},
		{
			[]byte{1, 1},
			int32(0),
		},
		{
			[]byte{1, 0, 1},
			int32(0),
		},
		{
			[]byte{1, 1, 1},
			int32(0),
		},
		// int64
		{
			nil,
			int64(0),
		},
		{
			[]byte{},
			int64(0),
		},
		{
			[]byte{1},
			int64(0),
		},
		{
			[]byte{1, 1},
			int64(0),
		},
		{
			[]byte{1, 1, 1},
			int64(0),
		},
		{
			[]byte{1, 1, 1, 1},
			int64(0),
		},
		{
			[]byte{1, 1, 1, 1, 1},
			int64(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1},
			int64(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1, 1},
			int64(0),
		},
		// int
		{
			nil,
			int(0),
		},
		{
			[]byte{},
			int(0),
		},
		{
			[]byte{1},
			int(0),
		},
		{
			[]byte{1, 1},
			int(0),
		},
		{
			[]byte{1, 1, 1},
			int(0),
		},
		{
			[]byte{1, 1, 1, 1},
			int(0),
		},
		{
			[]byte{1, 1, 1, 1, 1},
			int(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1},
			int(0),
		},
		{
			[]byte{1, 1, 1, 1, 1, 1, 1},
			int(0),
		},
	}
	for idx, testCase := range testCases {
		// allocate new value
		val := reflect.New(reflect.TypeOf(testCase.OutputValue))

		// unmarshal (valid values only)
		err := Unmarshal(testCase.Input, val.Interface())
		if err == nil {
			t.Errorf("expected Unmarshal error for testCase %d (%v), into %s/%s, but none was received",
				idx, testCase.Input, val.Type().String(), val.Kind().String())
			continue
		}
	}
}

func TestEncodeSlices(t *testing.T) {
	// bools
	boolSlice := []bool{false, true, false}
	b := Marshal(boolSlice)
	if len(b) != 4 {
		t.Errorf("expected bool slice to be of length 4, but it is %v instead", len(b))
	}
	var boolSliceCopy []bool
	err := Unmarshal(b, &boolSliceCopy)
	if err != nil {
		t.Error("error while unmarshalling bool slice", err)
	}
	if !reflect.DeepEqual(boolSlice, boolSliceCopy) {
		t.Error(boolSlice, "!=", boolSliceCopy)
	}

	// strings
	strSlice := []string{"one", "hello", ""}
	b = Marshal(strSlice)
	if len(b) != 12 {
		t.Errorf("expected string slice to be of length 12, but it is %v instead", len(b))
	}
	var strSliceCopy []string
	err = Unmarshal(b, &strSliceCopy)
	if err != nil {
		t.Error("error while unmarshalling string slice", err)
	}
	if !reflect.DeepEqual(strSlice, strSliceCopy) {
		t.Error(strSlice, "!=", strSliceCopy)
	}
}

func TestEncodeArray(t *testing.T) {
	// bools
	i64Array := [2]i64{4, 2}
	b := Marshal(i64Array)
	if len(b) != 16 {
		t.Errorf("expected i64 array to be of length 16, but it is %v instead", len(b))
	}
	var i64ArrayCopy [2]i64
	err := Unmarshal(b, &i64ArrayCopy)
	if err != nil {
		t.Error("error while unmarshalling i64 array", err)
	}
	if !reflect.DeepEqual(i64Array, i64ArrayCopy) {
		t.Error(i64Array, "!=", i64ArrayCopy)
	}
}

func TestEncodeByteArray(t *testing.T) {
	// bytes
	bArray := [2]byte{4, 2}
	b := Marshal(bArray)
	if len(b) != 2 {
		t.Errorf("expected bytes array to be of length 2, but it is %v instead", len(b))
	}
	var bArrayCopy [2]byte
	err := Unmarshal(b, &bArrayCopy)
	if err != nil {
		t.Error("error while unmarshalling bytes array", err)
	}
	if !reflect.DeepEqual(bArray, bArrayCopy) {
		t.Error(bArray, "!=", bArrayCopy)
	}
}

func TestEncodeStruct(t *testing.T) {
	type Foo struct {
		Bar    string
		Answer u8
	}
	foo := Foo{"bar", 42}
	b := Marshal(foo)
	if len(b) != 5 {
		t.Errorf("expected foo struct to be of length 5, but it is %v instead", len(b))
	}
	var fooCopy Foo
	err := Unmarshal(b, &fooCopy)
	if err != nil {
		t.Error("error while unmarshalling Foo struct", err)
	}
	if !reflect.DeepEqual(foo, fooCopy) {
		t.Error(foo, "!=", fooCopy)
	}
}

func TestEncodeMarshallAll(t *testing.T) {
	b := MarshalAll("bar", u8(42))
	if len(b) != 5 {
		t.Errorf("expected encoded values to be of length 5, but it is %v instead", len(b))
	}
	var foo struct {
		Bar    string
		Answer u8
	}
	err := Unmarshal(b, &foo)
	if err != nil {
		t.Error("error while unmarshalling Foo struct", err)
	}
	if foo.Bar != "bar" {
		t.Error(foo.Bar, "!=", "bar")
	}
	if foo.Answer != u8(42) {
		t.Error(foo.Answer, "!=", u8(42))
	}
}

func TestEncodeMap(t *testing.T) {
	var (
		b   []byte
		err interface{}
	)
	func() {
		defer func() {
			err = recover()
		}()
		b = Marshal(map[string]interface{}{"nil": nil})
	}()
	if err == nil {
		t.Fatal("managed to encode map as ", b)
	}
}
