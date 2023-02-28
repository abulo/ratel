package i18n

import (
	"reflect"
	"testing"

	"github.com/gookit/ini/v2"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestT(t *testing.T) {
	type args struct {
		lang string
		key  string
		args []interface{}
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
			if got := T(tt.args.lang, tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("T() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTr(t *testing.T) {
	type args struct {
		lang string
		key  string
		args []interface{}
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
			if got := Tr(tt.args.lang, tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("Tr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDt(t *testing.T) {
	type args struct {
		key  string
		args []interface{}
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
			if got := Dt(tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("Dt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefTr(t *testing.T) {
	type args struct {
		key  string
		args []interface{}
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
			if got := DefTr(tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("DefTr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		langDir   string
		defLang   string
		languages map[string]string
	}
	tests := []struct {
		name string
		args args
		want *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.langDir, tt.args.defLang, tt.args.languages); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		langDir   string
		defLang   string
		languages map[string]string
	}
	tests := []struct {
		name string
		args args
		want *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.langDir, tt.args.defLang, tt.args.languages); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEmpty(t *testing.T) {
	tests := []struct {
		name string
		want *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmpty(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithInit(t *testing.T) {
	type args struct {
		langDir   string
		defLang   string
		languages map[string]string
	}
	tests := []struct {
		name string
		args args
		want *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithInit(tt.args.langDir, tt.args.defLang, tt.args.languages); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_Dt(t *testing.T) {
	type args struct {
		key  string
		args []interface{}
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Dt(tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("I18n.Dt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_DefTr(t *testing.T) {
	type args struct {
		key  string
		args []interface{}
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.DefTr(tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("I18n.DefTr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_T(t *testing.T) {
	type args struct {
		lang string
		key  string
		args []interface{}
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.T(tt.args.lang, tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("I18n.T() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_Tr(t *testing.T) {
	type args struct {
		lang string
		key  string
		args []interface{}
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Tr(tt.args.lang, tt.args.key, tt.args.args...); got != tt.want {
				t.Errorf("I18n.Tr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_HasKey(t *testing.T) {
	type args struct {
		lang string
		key  string
	}
	tests := []struct {
		name   string
		l      *I18n
		args   args
		wantOk bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := tt.l.HasKey(tt.args.lang, tt.args.key); gotOk != tt.wantOk {
				t.Errorf("I18n.HasKey() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestI18n_transFromFallback(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.transFromFallback(tt.args.key); got != tt.want {
				t.Errorf("I18n.transFromFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_renderMessage(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.renderMessage(tt.args.msg, tt.args.args...); got != tt.want {
				t.Errorf("I18n.renderMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStr := toString(tt.args.val); gotStr != tt.wantStr {
				t.Errorf("toString() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func TestI18n_Init(t *testing.T) {
	tests := []struct {
		name string
		l    *I18n
		want *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("I18n.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_loadSingleFiles(t *testing.T) {
	tests := []struct {
		name string
		l    *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.loadSingleFiles()
		})
	}
}

func TestI18n_loadDirFiles(t *testing.T) {
	tests := []struct {
		name string
		l    *I18n
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.loadDirFiles()
		})
	}
}

func TestI18n_Add(t *testing.T) {
	type args struct {
		lang string
		name string
	}
	tests := []struct {
		name string
		l    *I18n
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Add(tt.args.lang, tt.args.name)
		})
	}
}

func TestI18n_NewLang(t *testing.T) {
	type args struct {
		lang string
		name string
	}
	tests := []struct {
		name string
		l    *I18n
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.NewLang(tt.args.lang, tt.args.name)
		})
	}
}

func TestI18n_LoadFile(t *testing.T) {
	type args struct {
		lang string
		file string
	}
	tests := []struct {
		name    string
		l       *I18n
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.LoadFile(tt.args.lang, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("I18n.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestI18n_LoadString(t *testing.T) {
	type args struct {
		lang string
		data string
	}
	tests := []struct {
		name    string
		l       *I18n
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.LoadString(tt.args.lang, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("I18n.LoadString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestI18n_Lang(t *testing.T) {
	type args struct {
		lang string
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want *ini.Ini
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Lang(tt.args.lang); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("I18n.Lang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_Export(t *testing.T) {
	type args struct {
		lang string
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Export(tt.args.lang); got != tt.want {
				t.Errorf("I18n.Export() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_HasLang(t *testing.T) {
	type args struct {
		lang string
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.HasLang(tt.args.lang); got != tt.want {
				t.Errorf("I18n.HasLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_DelLang(t *testing.T) {
	type args struct {
		lang string
	}
	tests := []struct {
		name string
		l    *I18n
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.DelLang(tt.args.lang); got != tt.want {
				t.Errorf("I18n.DelLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18n_Languages(t *testing.T) {
	tests := []struct {
		name string
		l    *I18n
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Languages(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("I18n.Languages() = %v, want %v", got, tt.want)
			}
		})
	}
}
