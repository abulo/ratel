package sql

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func Test_postgresOpen(t *testing.T) {
	type args struct {
		driverName string
		dns        string
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := postgresOpen(tt.args.driverName, tt.args.dns)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgresOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgresOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}
