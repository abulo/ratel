package util

import (
	"html/template"
	"reflect"
	"testing"
	_ "time/tzdata"
)

func TestMarshalHTML(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want template.HTML
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarshalHTML(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalJS(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want template.JS
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarshalJS(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJS(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want template.JS
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JS(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatic(t *testing.T) {
	type args struct {
		p string
		v string
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
			if got := Static(tt.args.p, tt.args.v); got != tt.want {
				t.Errorf("Static() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainStatic(t *testing.T) {
	type args struct {
		p string
		v string
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
			if got := DomainStatic(tt.args.p, tt.args.v); got != tt.want {
				t.Errorf("DomainStatic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnescapeString(t *testing.T) {
	type args struct {
		v any
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
			if got := UnescapeString(tt.args.v); got != tt.want {
				t.Errorf("UnescapeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppRootPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppRootPath(); got != tt.want {
				t.Errorf("GetAppRootPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_substr(t *testing.T) {
	type args struct {
		s      string
		pos    int
		length int
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
			if got := substr(tt.args.s, tt.args.pos, tt.args.length); got != tt.want {
				t.Errorf("substr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetParentDirectory(t *testing.T) {
	type args struct {
		dir string
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
			if got := GetParentDirectory(tt.args.dir); got != tt.want {
				t.Errorf("GetParentDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentDirectory(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCurrentDirectory(); got != tt.want {
				t.Errorf("GetCurrentDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandom(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Random(); got != tt.want {
				t.Errorf("Random() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyPhone(t *testing.T) {
	type args struct {
		phone string
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
			if got := VerifyPhone(tt.args.phone); got != tt.want {
				t.Errorf("VerifyPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	type args struct {
		email string
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
			if got := VerifyEmail(tt.args.email); got != tt.want {
				t.Errorf("VerifyEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyIPv4(t *testing.T) {
	type args struct {
		address string
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
			if got := VerifyIPv4(tt.args.address); got != tt.want {
				t.Errorf("VerifyIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyIPv6(t *testing.T) {
	type args struct {
		address string
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
			if got := VerifyIPv6(tt.args.address); got != tt.want {
				t.Errorf("VerifyIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionName(t *testing.T) {
	type args struct {
		i any
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
			if got := FunctionName(tt.args.i); got != tt.want {
				t.Errorf("FunctionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intSlice_Len(t *testing.T) {
	tests := []struct {
		name string
		s    intSlice
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("intSlice.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intSlice_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    intSlice
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Swap(tt.args.i, tt.args.j)
		})
	}
}

func Test_intSlice_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    intSlice
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("intSlice.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_equal(t *testing.T) {
	type args struct {
		a []int
		b []int
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
			if got := equal(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMissingElement(t *testing.T) {
	type args struct {
		arr []int
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
			if got := GetMissingElement(tt.args.arr); got != tt.want {
				t.Errorf("GetMissingElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReplacer(t *testing.T) {
	type args struct {
		endpoint string
		values   []any
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
			if got := NewReplacer(tt.args.endpoint, tt.args.values...); got != tt.want {
				t.Errorf("NewReplacer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteReplace(t *testing.T) {
	type args struct {
		s   []byte
		old []byte
		new []byte
		n   int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteReplace(tt.args.s, tt.args.old, tt.args.new, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByteReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}
