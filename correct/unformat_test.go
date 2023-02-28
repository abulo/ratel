package correct

import "testing"

func TestUnformat(t *testing.T) {
	type args struct {
		text    string
		options []UnformatOption
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
			if got := Unformat(tt.args.text, tt.args.options...); got != tt.want {
				t.Errorf("Unformat() = %v, want %v", got, tt.want)
			}
		})
	}
}
