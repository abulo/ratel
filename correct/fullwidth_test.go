package correct

import "testing"

func Test_fullwidth(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := fullwidth(tt.args.text); gotOut != tt.wantOut {
				t.Errorf("fullwidth() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_fullwidthReplacePart(t *testing.T) {
	type args struct {
		part string
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
			if got := fullwidthReplacePart(tt.args.part); got != tt.want {
				t.Errorf("fullwidthReplacePart() = %v, want %v", got, tt.want)
			}
		})
	}
}
