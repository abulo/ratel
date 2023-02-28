package correct

import (
	"reflect"
	"testing"
)

func Test_newStrategery(t *testing.T) {
	type args struct {
		one     string
		other   string
		space   bool
		reverse bool
	}
	tests := []struct {
		name string
		args args
		want *strategery
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newStrategery(tt.args.one, tt.args.other, tt.args.space, tt.args.reverse); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newStrategery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strategery_format(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		s       *strategery
		args    args
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := tt.s.format(tt.args.in); gotOut != tt.wantOut {
				t.Errorf("strategery.format() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_strategery_addSpace(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		s       *strategery
		args    args
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := tt.s.addSpace(tt.args.in); gotOut != tt.wantOut {
				t.Errorf("strategery.addSpace() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_strategery_removeSpace(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		s       *strategery
		args    args
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := tt.s.removeSpace(tt.args.in); gotOut != tt.wantOut {
				t.Errorf("strategery.removeSpace() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
