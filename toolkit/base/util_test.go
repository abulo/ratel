package base

import "testing"

func TestCamelStr(t *testing.T) {
	type args struct {
		name string
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
			if got := CamelStr(tt.args.name); got != tt.want {
				t.Errorf("CamelStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelper(t *testing.T) {
	type args struct {
		name string
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
			if got := Helper(tt.args.name); got != tt.want {
				t.Errorf("Helper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChar(t *testing.T) {
	type args struct {
		in string
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
			if got := Char(tt.args.in); got != tt.want {
				t.Errorf("Char() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSymbolChar(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SymbolChar(); got != tt.want {
				t.Errorf("SymbolChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		numberOne any
		numberTwo any
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
			if got := Add(tt.args.numberOne, tt.args.numberTwo); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	type args struct {
		Condition []Column
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
			if got := Convert(tt.args.Condition); got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleDaoConvertProto(t *testing.T) {
	type args struct {
		Condition []Column
		res       string
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
			if got := ModuleDaoConvertProto(tt.args.Condition, tt.args.res); got != tt.want {
				t.Errorf("ModuleDaoConvertProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiToProto(t *testing.T) {
	type args struct {
		Condition []Column
		res       string
		request   string
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
			if got := ApiToProto(tt.args.Condition, tt.args.res, tt.args.request); got != tt.want {
				t.Errorf("ApiToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleProtoConvertMap(t *testing.T) {
	type args struct {
		Condition []Column
		request   string
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
			if got := ModuleProtoConvertMap(tt.args.Condition, tt.args.request); got != tt.want {
				t.Errorf("ModuleProtoConvertMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleProtoConvertDao(t *testing.T) {
	type args struct {
		Condition []Column
		res       string
		request   string
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
			if got := ModuleProtoConvertDao(tt.args.Condition, tt.args.res, tt.args.request); got != tt.want {
				t.Errorf("ModuleProtoConvertDao() = %v, want %v", got, tt.want)
			}
		})
	}
}
