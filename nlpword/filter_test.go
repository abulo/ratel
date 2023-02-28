package nlpword

import (
	"io"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Filter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_UpdateNoisePattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.filter.UpdateNoisePattern(tt.args.pattern)
		})
	}
}

func TestFilter_LoadWordDict(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		filter  *Filter
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.filter.LoadWordDict(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Filter.LoadWordDict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilter_LoadNetWordDict(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		filter  *Filter
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.filter.LoadNetWordDict(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Filter.LoadNetWordDict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilter_Load(t *testing.T) {
	type args struct {
		rd io.Reader
	}
	tests := []struct {
		name    string
		filter  *Filter
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.filter.Load(tt.args.rd); (err != nil) != tt.wantErr {
				t.Errorf("Filter.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilter_updateFailureLink(t *testing.T) {
	tests := []struct {
		name   string
		filter *Filter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.filter.updateFailureLink()
		})
	}
}

func TestFilter_AddWord(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.filter.AddWord(tt.args.words...)
		})
	}
}

func TestFilter_Filter(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.Filter(tt.args.text); got != tt.want {
				t.Errorf("Filter.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Replace(t *testing.T) {
	type args struct {
		text string
		repl rune
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.Replace(tt.args.text, tt.args.repl); got != tt.want {
				t.Errorf("Filter.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_FindIn(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
		want   bool
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.filter.FindIn(tt.args.text)
			if got != tt.want {
				t.Errorf("Filter.FindIn() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Filter.FindIn() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFilter_FindAll(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.FindAll(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Validate(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
		want   bool
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.filter.Validate(tt.args.text)
			if got != tt.want {
				t.Errorf("Filter.Validate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Filter.Validate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFilter_RemoveNoise(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name   string
		filter *Filter
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.RemoveNoise(tt.args.text); got != tt.want {
				t.Errorf("Filter.RemoveNoise() = %v, want %v", got, tt.want)
			}
		})
	}
}
