package util

import (
	"reflect"
	"testing"
)

func TestArrayFill(t *testing.T) {
	type args struct {
		startIndex int
		num        uint
		value      any
	}
	tests := []struct {
		name string
		args args
		want map[int]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayFill(tt.args.startIndex, tt.args.num, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayFill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayFlip(t *testing.T) {
	type args struct {
		m map[any]any
	}
	tests := []struct {
		name string
		args args
		want map[any]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayFlip(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayFlip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayKeys(t *testing.T) {
	type args struct {
		elements map[any]any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayKeys(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayValues(t *testing.T) {
	type args struct {
		elements map[any]any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayValues(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayMerge(t *testing.T) {
	type args struct {
		ss [][]any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayMerge(tt.args.ss...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayChunk(t *testing.T) {
	type args struct {
		s    []any
		size int
	}
	tests := []struct {
		name string
		args args
		want [][]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayChunk(tt.args.s, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayChunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayPad(t *testing.T) {
	type args struct {
		s    []any
		size int
		val  any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayPad(tt.args.s, tt.args.size, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayPad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArraySlice(t *testing.T) {
	type args struct {
		s      []any
		offset uint
		length uint
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArraySlice(tt.args.s, tt.args.offset, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArraySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayRand(t *testing.T) {
	type args struct {
		elements []any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayRand(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayRand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayColumn(t *testing.T) {
	type args struct {
		input     map[string]map[string]any
		columnKey string
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayColumn(tt.args.input, tt.args.columnKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayPush(t *testing.T) {
	type args struct {
		s        *[]any
		elements []any
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayPush(tt.args.s, tt.args.elements...); got != tt.want {
				t.Errorf("ArrayPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayPop(t *testing.T) {
	type args struct {
		s *[]any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayPop(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayPop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayUnshift(t *testing.T) {
	type args struct {
		s        *[]any
		elements []any
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayUnshift(tt.args.s, tt.args.elements...); got != tt.want {
				t.Errorf("ArrayUnshift() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayShift(t *testing.T) {
	type args struct {
		s *[]any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayShift(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayShift() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayKeyExists(t *testing.T) {
	type args struct {
		key any
		m   map[any]any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayKeyExists(tt.args.key, tt.args.m); got != tt.want {
				t.Errorf("ArrayKeyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayCombine(t *testing.T) {
	type args struct {
		s1 []any
		s2 []any
	}
	tests := []struct {
		name string
		args args
		want map[any]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayCombine(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayCombine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayReverse(t *testing.T) {
	type args struct {
		s []any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayReverse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayReverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImplode(t *testing.T) {
	type args struct {
		glue   string
		pieces []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Implode(tt.args.glue, tt.args.pieces); got != tt.want {
				t.Errorf("Implode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInArray(t *testing.T) {
	type args struct {
		needle   any
		haystack any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InArray(tt.args.needle, tt.args.haystack); got != tt.want {
				t.Errorf("InArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayRandMap(t *testing.T) {
	type args struct {
		elements []map[string]string
	}
	tests := []struct {
		name string
		args args
		want []map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayRandMap(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayRandMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMultiArray(t *testing.T) {
	type args struct {
		haystack any
		needle   []any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InMultiArray(tt.args.haystack, tt.args.needle...); got != tt.want {
				t.Errorf("InMultiArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiArray(t *testing.T) {
	type args struct {
		haystack any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MultiArray(tt.args.haystack); got != tt.want {
				t.Errorf("MultiArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceToString(t *testing.T) {
	type args struct {
		data []any
	}
	tests := []struct {
		name  string
		args  args
		wantS []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := InterfaceToString(tt.args.data); !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("InterfaceToString() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestSplitString(t *testing.T) {
	type args struct {
		p     string
		split string
		space bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitString(tt.args.p, tt.args.split, tt.args.space); got != tt.want {
				t.Errorf("SplitString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayPluck(t *testing.T) {
	type args struct {
		data  []map[string]string
		value string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayPluck(tt.args.data, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayPluck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayRemoveRepeatedElement(t *testing.T) {
	type args struct {
		arr []string
	}
	tests := []struct {
		name       string
		args       args
		wantNewArr []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewArr := ArrayRemoveRepeatedElement(tt.args.arr); !reflect.DeepEqual(gotNewArr, tt.wantNewArr) {
				t.Errorf("ArrayRemoveRepeatedElement() = %v, want %v", gotNewArr, tt.wantNewArr)
			}
		})
	}
}

func TestArrayKeyPluck(t *testing.T) {
	type args struct {
		data  []map[string]string
		value string
		key   string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayKeyPluck(tt.args.data, tt.args.value, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayKeyPluck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayMultiPluck(t *testing.T) {
	type args struct {
		data []map[string]string
		key  string
	}
	tests := []struct {
		name string
		args args
		want map[string]map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayMultiPluck(tt.args.data, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayMultiPluck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonStringToObject(t *testing.T) {
	type args struct {
		s string
		v any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsonStringToObject(tt.args.s, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("jsonStringToObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAryMapStringToAryMapInterface(t *testing.T) {
	type args struct {
		d []map[string]string
	}
	tests := []struct {
		name string
		args args
		want []map[string]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AryMapStringToAryMapInterface(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AryMapStringToAryMapInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapStringToMapInterface(t *testing.T) {
	type args struct {
		d map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapStringToMapInterface(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapStringToMapInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAryMapInterfaceToAryMapString(t *testing.T) {
	type args struct {
		d []map[string]any
	}
	tests := []struct {
		name string
		args args
		want []map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AryMapInterfaceToAryMapString(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AryMapInterfaceToAryMapString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapInterfaceToMapString(t *testing.T) {
	type args struct {
		d map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapInterfaceToMapString(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapInterfaceToMapString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgStringToAryInterface(t *testing.T) {
	type args struct {
		d []string
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArgStringToAryInterface(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArgStringToAryInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAryInterfaceToArgString(t *testing.T) {
	type args struct {
		d []any
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AryInterfaceToArgString(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AryInterfaceToArgString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceToAryMapStringInterface(t *testing.T) {
	type args struct {
		in any
	}
	tests := []struct {
		name string
		args args
		want []map[string]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceToAryMapStringInterface(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceToAryMapStringInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceToAryMapStringString(t *testing.T) {
	type args struct {
		in any
	}
	tests := []struct {
		name string
		args args
		want []map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceToAryMapStringString(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceToAryMapStringString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayStringUniq(t *testing.T) {
	type args struct {
		arr []string
	}
	tests := []struct {
		name       string
		args       args
		wantNewArr []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewArr := ArrayStringUniq(tt.args.arr); !reflect.DeepEqual(gotNewArr, tt.wantNewArr) {
				t.Errorf("ArrayStringUniq() = %v, want %v", gotNewArr, tt.wantNewArr)
			}
		})
	}
}
