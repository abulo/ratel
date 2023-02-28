package correct

import "testing"

func Test_registerStrategery(t *testing.T) {
	type args struct {
		one     string
		other   string
		space   bool
		reverse bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerStrategery(tt.args.one, tt.args.other, tt.args.space, tt.args.reverse)
		})
	}
}

func Test_spaceDashWithHans(t *testing.T) {
	type args struct {
		in string
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
			if gotOut := spaceDashWithHans(tt.args.in); gotOut != tt.wantOut {
				t.Errorf("spaceDashWithHans() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type args struct {
		in      string
		options []Option
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
			if gotOut := Format(tt.args.in, tt.args.options...); gotOut != tt.wantOut {
				t.Errorf("Format() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
