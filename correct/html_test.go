package correct

import "testing"

func TestFormatHTML(t *testing.T) {
	type args struct {
		body    string
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := FormatHTML(tt.args.body, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("FormatHTML() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestUnformatHTML(t *testing.T) {
	type args struct {
		body    string
		options []UnformatOption
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := UnformatHTML(tt.args.body, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnformatHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("UnformatHTML() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_processHTML(t *testing.T) {
	type args struct {
		body string
		fn   func(plainText string) string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := processHTML(tt.args.body, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("processHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("processHTML() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
