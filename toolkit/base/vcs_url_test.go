package base

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseVCSUrl(t *testing.T) {
	type args struct {
		repo string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVCSUrl(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVCSUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVCSUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
